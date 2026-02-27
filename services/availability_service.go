package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type prometheusInstantResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			// Value format: [ <unix_time>, "<string_value>" ]
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
	Error     string `json:"error,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
}

// StartSystemAvailabilityJob runs the availability snapshot daily at 09:59 UTC
// and writes it into transactions.system_availability.
func StartSystemAvailabilityJob() {
	go func() {
		for {
			now := time.Now().UTC()
			nextRun := time.Date(now.Year(), now.Month(), now.Day(), 9, 59, 0, 0, time.UTC)
			if !now.Before(nextRun) {
				nextRun = nextRun.Add(24 * time.Hour)
			}

			timer := time.NewTimer(nextRun.Sub(now))
			<-timer.C

			runAvailabilitySnapshot()
		}
	}()
}

// runAvailabilitySnapshot pulls metrics and persists one row. Separated for easier testing.
func runAvailabilitySnapshot() {
	cfg := model.LoadDBConfig()
	if err := config.ValidateDBConfig(cfg); err != nil {
		log.Printf("[ERROR] availability job skipped (db config): %v", err)
		return
	}

	rangeStr := "24h"

	success, err := queryPrometheusValue(fmt.Sprintf(`sum(increase(nginx_ingress_controller_requests{status=~"(2|4).."}[%s]))`, rangeStr))
	if err != nil {
		log.Printf("[ERROR] availability job: fetch success metric: %v", err)
		return
	}

	total, err := queryPrometheusValue(fmt.Sprintf(`sum(increase(nginx_ingress_controller_requests[%s]))`, rangeStr))
	if err != nil {
		log.Printf("[ERROR] availability job: fetch total metric: %v", err)
		return
	}

	errorReq := total - success
	if errorReq < 0 {
		errorReq = 0
	}

	var availability float64
	if total > 0 {
		availability = (success / total) * 100
		availability = math.Round(availability*1000) / 1000 // keep 3 decimal places
	} else {
		availability = 0
	}

	db, err := model.ConnectTransDB(cfg)
	if err != nil {
		log.Printf("[ERROR] availability job: connect trans db: %v", err)
		return
	}
	defer db.Close()

	row := model.SystemAvailability{
		TS:                  time.Now().UTC(),
		WindowInterval:      rangeStr,
		TotalRequests:       int64(math.Round(total)),
		SuccessRequests:     int64(math.Round(success)),
		ErrorRequests:       int64(math.Round(errorReq)),
		AvailabilityPercent: availability,
		P95LatencyMs:        nil,
		P99LatencyMs:        nil,
	}

	if err := db.InsertSystemAvailability(row); err != nil {
		log.Printf("[ERROR] availability job: insert row: %v", err)
	}
}

// queryPrometheusValue executes a Prometheus instant query and returns the numeric value.
func queryPrometheusValue(query string) (float64, error) {
	encoded := url.QueryEscape(query)
	reqURL := fmt.Sprintf("%s?query=%s", config.PrometheusUrl, encoded)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var pr prometheusInstantResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return 0, err
	}

	if pr.Status != "success" || len(pr.Data.Result) == 0 || len(pr.Data.Result[0].Value) < 2 {
		return 0, fmt.Errorf("unexpected prometheus response: status=%s len(result)=%d", pr.Status, len(pr.Data.Result))
	}

	valueStr, ok := pr.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("cannot parse value field")
	}

	return strconv.ParseFloat(valueStr, 64)
}
