package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/muchrm/go-healthcheck/config"
	"github.com/muchrm/go-healthcheck/healthcheck"
	"github.com/muchrm/go-healthcheck/lineapi"
	"github.com/muchrm/go-healthcheck/utils"
)

func main() {
	config.InitConfig()
	CSVPath, err := config.GetString(config.CSVPath)
	if err != nil {
		fmt.Println(config.CSVPathMissingError)
		os.Exit(0)
	}
	timeout, err := config.GetInt(config.HealthcheckTimeoutInSecond)
	if err != nil {
		timeout = 30
	}

	clientWithoutTimeout := &http.Client{}
	clientWithTimeout := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	URLS, err := utils.ReadCSVAndGetServerUrls(CSVPath)
	healthCheckResult := healthcheck.RunHealthCheck(clientWithTimeout, URLS)
	healthcheck.PrintHealthCheckSummary(healthCheckResult)
	lineapi.SentResult(clientWithoutTimeout, healthCheckResult)

}
