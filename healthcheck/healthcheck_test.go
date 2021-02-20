package healthcheck

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestIsServerTimeout(t *testing.T) {
	type MockResponse struct {
		status int
		delay  time.Duration
		body   string
	}
	tests := []struct {
		name string
		mock MockResponse
		want bool
	}{
		{
			name: "HTTP Status Ok should not be timeout",
			mock: MockResponse{
				status: http.StatusOK,
				delay:  0,
				body:   "HealthCheck ok",
			},
			want: false,
		}, {
			name: "HTTP Status BadRequest not be timeout",
			mock: MockResponse{
				status: http.StatusBadRequest,
				delay:  0,
				body:   "HealthCheck ok",
			},
			want: false,
		}, {
			name: "10 second should be timeout",
			mock: MockResponse{
				status: http.StatusOK,
				delay:  10 * time.Second,
				body:   "HealthCheck not ok",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.mock.delay > 0 {
					time.Sleep(tt.mock.delay)
				}
				w.WriteHeader(tt.mock.status)
				fmt.Fprintln(w, tt.mock.body)
			}))
			defer ts.Close()

			client := &http.Client{
				Timeout: 10 * time.Millisecond,
			}
			if got := IsServerTimeout(client, ts.URL); got != tt.want {
				t.Errorf("getHealcheckStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunHealthCheck(t *testing.T) {
	tests := []struct {
		name                string
		numberOfSuccessUrls int
		numberOfFailureUrls int
		want                HealthCheckResult
	}{
		{
			name:                "set success url array 5 item should return HealthCheckResult correctly",
			numberOfSuccessUrls: 5,
			numberOfFailureUrls: 0,
			want: HealthCheckResult{
				TotalWebistes:   5,
				TotalSuccessful: 5,
				TotalFailure:    0,
			},
		},
		{
			name:                "set fail url array 5 item should return HealthCheckResult correctly",
			numberOfSuccessUrls: 0,
			numberOfFailureUrls: 5,
			want: HealthCheckResult{
				TotalWebistes:   5,
				TotalSuccessful: 0,
				TotalFailure:    5,
			},
		},
		{
			name:                "set fail and success url array 7 item should return HealthCheckResult correctly",
			numberOfSuccessUrls: 3,
			numberOfFailureUrls: 4,
			want: HealthCheckResult{
				TotalWebistes:   7,
				TotalSuccessful: 3,
				TotalFailure:    4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "ok")
			}))
			defer ts.Close()

			tsWithTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(10 * time.Second)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "ok")
			}))
			defer tsWithTimeout.Close()

			client := &http.Client{
				Timeout: time.Duration(10) * time.Millisecond,
			}
			checkURLs := []string{}
			for i := 0; i < tt.numberOfSuccessUrls; i++ {
				checkURLs = append(checkURLs, ts.URL)
			}
			for i := 0; i < tt.numberOfFailureUrls; i++ {
				checkURLs = append(checkURLs, tsWithTimeout.URL)
			}
			got := RunHealthCheck(client, checkURLs)
			got.Duration = 0
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RunHealthCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
