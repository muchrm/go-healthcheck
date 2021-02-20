package healthcheck

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/muchrm/go-healthcheck/config"
	"github.com/muchrm/go-healthcheck/utils"
)

type HealthCheckResult struct {
	TotalWebistes   int           `json:"total_websites"`
	TotalSuccessful int           `json:"success"`
	TotalFailure    int           `json:"failure"`
	Duration        time.Duration `json:"total_time"`
}

// IsServerTimeout receive url and return true only when http timeout
func IsServerTimeout(client *http.Client, url string) bool {
	_, err := client.Get(url)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return true
		}
	}
	return false
}

// RunHealthCheck receive url array and run check timeout
func RunHealthCheck(client *http.Client, URLS []string) HealthCheckResult {
	fmt.Println("Perform website checking...")

	start := time.Now()
	totalJobs := len(URLS)
	jobs := make(chan string, totalJobs)
	serverTimeoutResults := make(chan bool, totalJobs)

	maxWorker, err := config.GetInt(config.MaxWorker)
	if err != nil {
		maxWorker = 1
	}
	for w := 0; w < maxWorker; w++ {
		go createWorker(jobs, serverTimeoutResults, client)
	}

	for _, URLString := range URLS {
		jobs <- URLString
	}
	close(jobs)
	totalSuccessful := 0
	totalFailure := 0
	for range URLS {
		isTimeout := <-serverTimeoutResults
		if isTimeout {
			totalFailure++
		} else {
			totalSuccessful++
		}
	}

	elapsed := time.Since(start)
	return HealthCheckResult{
		TotalWebistes:   totalJobs,
		TotalSuccessful: totalSuccessful,
		TotalFailure:    totalFailure,
		Duration:        elapsed,
	}
}

// PrintHealthCheckSummary print healthcheck result in a human readable format
func PrintHealthCheckSummary(result HealthCheckResult) {
	fmt.Printf("Done!\n\n")
	fmt.Printf("Checked webistes: %d\n", result.TotalWebistes)
	fmt.Printf("Successful websites: %d\n", result.TotalSuccessful)
	fmt.Printf("Failure websites: %d \n", result.TotalFailure)
	fmt.Printf("Total times to finished checking website: %s", utils.GetDurationString((result.Duration)))
}

func createWorker(jobs <-chan string, results chan<- bool, client *http.Client) {
	for URLString := range jobs {
		results <- IsServerTimeout(client, URLString)
	}
}
