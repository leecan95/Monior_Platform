package services

import (
	"Monitor_Platform/model"
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	nginxIngressReadyQuery     = `kube_pod_status_ready{namespace="nginx-ingress", condition="true"}`
	nginxIngressReadyNamespace = "nginx-ingress"
)

// StartNginxIngressReadyStatusJob runs every 30s and stores nginx-ingress ready status.
func StartNginxIngressReadyStatusJob() {
	go func() {
		runNginxIngressReadySnapshot()

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			runNginxIngressReadySnapshot()
		}
	}()
}

func runNginxIngressReadySnapshot() {
	cfg := model.LoadDBConfig()

	pr, err := queryPrometheusInstantVector(nginxIngressReadyQuery)
	if err != nil {
		log.Printf("[ERROR] nginx ingress ready job: fetch prometheus vector: %v", err)
		return
	}

	readyStatus, promTS, err := calculateIngressReadyStatus(pr)
	if err != nil {
		log.Printf("[ERROR] nginx ingress ready job: parse vector result: %v", err)
		return
	}

	db, err := model.ConnectTransDB(cfg)
	if err != nil {
		log.Printf("[ERROR] nginx ingress ready job: connect trans db: %v", err)
		return
	}

	row := model.NginxIngressReadyStatus{
		CollectedAt:  time.Now().UTC(),
		PrometheusTS: promTS,
		Namespace:    nginxIngressReadyNamespace,
		ReadyStatus:  readyStatus,
	}

	if err := db.InsertNginxIngressReadyStatus(row); err != nil {
		log.Printf("[ERROR] nginx ingress ready job: insert row: %v", err)
	}
}

func calculateIngressReadyStatus(pr prometheusInstantResponse) (int16, *time.Time, error) {
	if len(pr.Data.Result) == 0 {
		return 0, nil, nil
	}

	readyStatus := int16(0)
	var latestPromTS *time.Time

	for _, result := range pr.Data.Result {
		if len(result.Value) < 2 {
			return 0, nil, fmt.Errorf("missing value field for pod=%s", result.Metric["pod"])
		}

		valueStr, ok := result.Value[1].(string)
		if !ok {
			return 0, nil, fmt.Errorf("invalid value type for pod=%s", result.Metric["pod"])
		}

		valueFloat, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return 0, nil, fmt.Errorf("parse readiness value for pod=%s: %w", result.Metric["pod"], err)
		}

		promTSRaw, ok := result.Value[0].(float64)
		if !ok {
			return 0, nil, fmt.Errorf("invalid timestamp type for pod=%s", result.Metric["pod"])
		}
		promTS := floatUnixToTime(promTSRaw)
		if latestPromTS == nil || promTS.After(*latestPromTS) {
			ts := promTS
			latestPromTS = &ts
		}

		if valueFloat > 0 {
			readyStatus = 1
		}
	}

	return readyStatus, latestPromTS, nil
}
