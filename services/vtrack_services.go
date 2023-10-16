package services

import (
	"Monitor_Platform/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetUsersTPS(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = sum(irate(users_api_request_count[30s]))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetOrganizationTPS(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = round(sum(irate(organizations_api_request_count[30s])))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetAdapterTPS(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = round(sum(irate(vtadapter_api_request_count{method=\"HandleKafka\"}[30s])) by (method), 0.001)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetRoleTPS(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetDeviceTPS(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = round(sum(irate(vtdevices_api_request_count[30s])))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetDataUsageStreaming(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = node_network_receive_bytes_total{device=\"br-b6703b71de00\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetUrlState(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = probe_success{}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetPodState(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = sum(kube_pod_status_phase{})"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetKafkaPartitionOnline(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = sum(kafka_server_replicamanager_partitioncount{job=~\"$job\"})"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetKafkaBytein(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = sum (rate(kafka_server_brokertopicmetrics_bytesin_total{}[1m]))by(instance)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetKafkaByteout(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = sum (rate(kafka_server_brokertopicmetrics_bytesout_total{}[1m]))by(instance)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}

func GetPrometheus(c *gin.Context, query string) interface{} {
	url := config.PrometheusUrl
	params := "?query=" + query
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp

}
