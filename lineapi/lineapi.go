package lineapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muchrm/go-healthcheck/config"
	"github.com/muchrm/go-healthcheck/healthcheck"
)

// SentResult receive healthcheck and sent to line healthcheck server
// It's will ignore if line api token is missing
func SentResult(client *http.Client, healthCheckResult healthcheck.HealthCheckResult) {
	lineAPIToken, err := config.GetString(config.AccessToken)
	if err != nil || lineAPIToken == "" {
		fmt.Println("Access Token not found ignore send report")
		return
	}
	body, _ := json.Marshal(healthCheckResult)
	url, _ := config.GetString(config.HealcheckReportAPI)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		fmt.Println("NewRequest Error: skipp send report")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lineAPIToken))
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("SentResult Error: Line Login Access Token not found skipp send report")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("SentResult Error: Unauthorized")
	} else if resp.StatusCode != http.StatusOK {
		fmt.Println("SentResult Error: Unhandled")
	}
}
