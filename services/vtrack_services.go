package services

import (
	"Monitor_Platform/config"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func GetUsersTPS(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=sum(irate(users_api_request_count[30s]))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}
	fmt.Printf("User reponse : %s", body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.JSON(500, err)
		return ""
	}
	value := result[0].(map[string]interface{})["value"].([]interface{})
	number := value[1].(string)
	return number
}

func GetOrganizationTPS(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=round(sum(irate(organizations_api_request_count[30s])))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.JSON(500, err)
		return ""
	}
	value := result[0].(map[string]interface{})["value"].([]interface{})
	number := value[1].(string)
	return number
}

func GetAdapterTPS(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=round(sum(irate(vtadapter_api_request_count{method=\"HandleKafka\"}[30s]))by(method))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.JSON(500, err)
		return ""
	}
	value := result[0].(map[string]interface{})["value"].([]interface{})
	number := value[1].(string)
	return number
}

func GetRoleTPS(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=round(sum(irate(mongodb_roles_request_count[3m])))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.JSON(500, err)
		return ""
	}
	value := result[0].(map[string]interface{})["value"].([]interface{})
	number := value[1].(string)
	return number
}

func GetDeviceTPS(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=round(sum(irate(vtdevices_api_request_count[3m])))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.JSON(500, err)
		return ""
	}
	value := result[0].(map[string]interface{})["value"].([]interface{})
	number := value[1].(string)
	return number
}

func GetDataUsageStreaming(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=node_network_receive_bytes_total{device=\"br-b6703b71de00\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}

	return response
}

func GetUrlState(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=probe_success{}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}

	return response
}

func GetPodState(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=sum(kube_pod_status_phase{})"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}

	return response
}

func GetKafkaPartitionOnline(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=sum(kafka_server_replicamanager_partitioncount{job=~\"$job\"})"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}

	return response
}

func GetKafkaBytein(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=sum(rate(kafka_server_brokertopicmetrics_bytesin_total{}[1m]))by(instance)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}

	return response
}

func GetKafkaByteout(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=sum(rate(kafka_server_brokertopicmetrics_bytesout_total{}[1m]))by(instance)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}

	return response
}

func GetPrometheus(c *gin.Context, query string) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=" + query
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.JSON(500, err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.JSON(500, err)
		return ""
	}

	return response

}
