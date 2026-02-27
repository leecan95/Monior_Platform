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

const (
	deploymentReplicaStatusQuery     = `kube_deployment_status_replicas_available{namespace="viot-app"}`
	deploymentReplicaStatusNamespace = "viot-app"
)

var monitoredDeployments = []string{
	"attributes",
	"devicegroups",
	"devices",
	"devicetokens",
	"executor",
	"geocoding",
	"imageserver",
	"miniovt",
	"organizations",
	"resolver",
	"roles",
	"users",
	"vtracking",
}

var monitoredDeploymentSet = func() map[string]struct{} {
	set := make(map[string]struct{}, len(monitoredDeployments))
	for _, deployment := range monitoredDeployments {
		set[deployment] = struct{}{}
	}
	return set
}()

// StartDeploymentReplicaStatusJob runs every 10s and stores replicas_available for configured deployments.
func StartDeploymentReplicaStatusJob() {
	go func() {
		runDeploymentReplicaStatusSnapshot()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			runDeploymentReplicaStatusSnapshot()
		}
	}()
}

func runDeploymentReplicaStatusSnapshot() {
	cfg := model.LoadDBConfig()
	if err := config.ValidateDBConfig(cfg); err != nil {
		log.Printf("[ERROR] deployment replica job skipped (db config): %v", err)
		return
	}

	pr, err := queryPrometheusInstantVector(deploymentReplicaStatusQuery)
	if err != nil {
		log.Printf("[ERROR] deployment replica job: fetch prometheus vector: %v", err)
		return
	}

	collectionTS := time.Now().UTC()
	rows, err := buildDeploymentReplicaRows(pr, collectionTS)
	if err != nil {
		log.Printf("[ERROR] deployment replica job: parse vector result: %v", err)
		return
	}

	db, err := model.ConnectTransDB(cfg)
	if err != nil {
		log.Printf("[ERROR] deployment replica job: connect trans db: %v", err)
		return
	}
	defer db.Close()

	if err := db.BatchInsertDeploymentReplicaStatuses(rows); err != nil {
		log.Printf("[ERROR] deployment replica job: insert rows: %v", err)
	}
}

func queryPrometheusInstantVector(query string) (prometheusInstantResponse, error) {
	encoded := url.QueryEscape(query)
	reqURL := fmt.Sprintf("%s?query=%s", config.PrometheusUrl, encoded)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return prometheusInstantResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return prometheusInstantResponse{}, err
	}
	defer resp.Body.Close()

	var pr prometheusInstantResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return prometheusInstantResponse{}, err
	}

	if pr.Status != "success" {
		return prometheusInstantResponse{}, fmt.Errorf("unexpected prometheus response: status=%s errorType=%s error=%s", pr.Status, pr.ErrorType, pr.Error)
	}

	return pr, nil
}

func buildDeploymentReplicaRows(pr prometheusInstantResponse, collectionTS time.Time) ([]model.DeploymentReplicaStatus, error) {
	replicasByDeployment := make(map[string]int64, len(monitoredDeployments))
	promTSByDeployment := make(map[string]time.Time, len(monitoredDeployments))

	for _, result := range pr.Data.Result {
		deployment := result.Metric["deployment"]
		if _, ok := monitoredDeploymentSet[deployment]; !ok {
			continue
		}
		if len(result.Value) < 2 {
			return nil, fmt.Errorf("missing value field for deployment=%s", deployment)
		}

		valueStr, ok := result.Value[1].(string)
		if !ok {
			return nil, fmt.Errorf("invalid value type for deployment=%s", deployment)
		}

		replicasFloat, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse replicas for deployment=%s: %w", deployment, err)
		}

		promTSRaw, ok := result.Value[0].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid timestamp type for deployment=%s", deployment)
		}

		replicasByDeployment[deployment] = int64(math.Round(replicasFloat))
		promTSByDeployment[deployment] = floatUnixToTime(promTSRaw)
	}

	rows := make([]model.DeploymentReplicaStatus, 0, len(monitoredDeployments))
	for _, deployment := range monitoredDeployments {
		replicas := int64(0)
		promTS := collectionTS

		if v, ok := replicasByDeployment[deployment]; ok {
			replicas = v
		} else {
			log.Printf("[WARN] deployment replica job: deployment=%s not found in prometheus response, default replicas=0", deployment)
		}

		if ts, ok := promTSByDeployment[deployment]; ok {
			promTS = ts
		}

		rows = append(rows, model.DeploymentReplicaStatus{
			CollectedAt:       collectionTS,
			PrometheusTS:      promTS,
			Namespace:         deploymentReplicaStatusNamespace,
			Deployment:        deployment,
			AvailableReplicas: replicas,
		})
	}

	return rows, nil
}

func floatUnixToTime(v float64) time.Time {
	sec, frac := math.Modf(v)
	nsec := int64(frac * float64(time.Second))
	return time.Unix(int64(sec), nsec).UTC()
}
