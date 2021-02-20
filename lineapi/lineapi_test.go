package lineapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/muchrm/go-healthcheck/config"
	"github.com/muchrm/go-healthcheck/healthcheck"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSentResult(t *testing.T) {
	t.Run("should run in correct format json", func(t *testing.T) {
		pass := false
		expected := []byte(`{"total_websites":7,"success":5,"failure":2,"total_time":4340000}`)
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost {
				got, _ := ioutil.ReadAll(r.Body)
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("TestSentResult() = %v, want %v", got, expected)
				} else {
					pass = true
				}
			}
			fmt.Fprintln(w, "ok")
		}))
		defer ts.Close()
		client := ts.Client()
		viper.Set(config.HealcheckReportAPI, ts.URL)
		viper.Set(config.AccessToken, "123456")
		SentResult(client, healthcheck.HealthCheckResult{
			TotalWebistes:   7,
			TotalSuccessful: 5,
			TotalFailure:    2,
			Duration:        4340000,
		})
		assert.Equal(t, true, pass)
	})

	t.Run("should not call http request when token not found", func(t *testing.T) {
		pass := false
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pass = true
			fmt.Fprintln(w, "ok")
		}))
		defer ts.Close()
		client := ts.Client()
		viper.Set(config.HealcheckReportAPI, ts.URL)
		viper.Set(config.AccessToken, nil)
		SentResult(client, healthcheck.HealthCheckResult{
			TotalWebistes:   7,
			TotalSuccessful: 5,
			TotalFailure:    2,
			Duration:        4340000,
		})
		assert.Equal(t, false, pass)
	})
}
