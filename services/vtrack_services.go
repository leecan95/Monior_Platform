package services

import (
	"Monitor_Platform/config"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

func GetUsersTPS(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=sum(irate(users_api_request_count[30s]))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	fmt.Printf("User reponse : %s", body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.Error(err)
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
		c.Error(err)
	}

	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.Error(err)
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
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.Error(err)
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
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.Error(err)
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
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	result := response["data"].(map[string]interface{})["result"].([]interface{})
	if result == nil {
		log.Printf("Wrong data for get number")
		c.Error(err)
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
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	return response
}

func GetUrlState(c *gin.Context) interface{} {
	var response KafkaReponse
	url := config.PrometheusUrl
	var sum []UrlData
	params := "?query=probe_success{system=~\"viettelmap\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := UrlData{
				Url:   result.Metric.Instance,
				Value: "Up",
			}
			sum = append(sum, data)
		} else {
			data := UrlData{
				Url:   result.Metric.Instance,
				Value: "Down",
			}
			sum = append(sum, data)
		}
	}

	return sum
}

func GetHttpStatusCode(c *gin.Context) interface{} {
	var response KafkaReponse
	url := config.PrometheusUrl
	var sum []UrlData
	params := "?query=probe_http_status_code{system=~\"viettelmap\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := UrlData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetPodState(c *gin.Context) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=sum(kube_pod_status_phase{})"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	return response
}

func GetKafkaPartitionOnline(c *gin.Context) interface{} {
	var response KafkaReponse
	url := config.PrometheusUrl
	params := "?query=sum(kafka_server_replicamanager_partitioncount{})"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	return response.Data.Result[0].Value[1]
}

func GetKafkaBytein(c *gin.Context) KafkaData {
	var response KafkaReponse
	url := config.PrometheusUrl
	params := "?query=sum(rate(kafka_server_brokertopicmetrics_bytesin_total{}[1m]))by(instance)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	return response.Data
}

func GetKafkaByteout(c *gin.Context) KafkaData {
	var response KafkaReponse
	url := config.PrometheusUrl
	params := "?query=sum(rate(kafka_server_brokertopicmetrics_bytesout_total{}[1m]))by(instance)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	return response.Data
}

func GetPodPend(c *gin.Context) []PodData {
	var response PodReponse
	var sum []PodData
	url := config.PrometheusUrl
	params := "?query=kube_pod_status_phase{pod=~\"^attributes.*\",phase=\"Pending\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}

	params = "?query=kube_pod_status_phase{pod=~\"^roles.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^organizations.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^users.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^vtracking.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}

	return sum
}

func GetPodUp(c *gin.Context) []PodData {
	var response PodReponse
	var sum []PodData
	url := config.PrometheusUrl
	params := "?query=kube_pod_status_phase{pod=~\"^attributes.*\",phase=\"Running\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}

	params = "?query=kube_pod_status_phase{pod=~\"^roles.*\",phase=\"Running\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^organizations.*\",phase=\"Running\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^users.*\",phase=\"Running\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^vtracking.*\",phase=\"Running\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^attributes.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "2",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^roles.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "2",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^organizations.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "2",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^users.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "2",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^vtracking.*\",phase=\"Pending\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "2",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^attributes.*\",phase=\"Failed\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "0",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^roles.*\",phase=\"Failed\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "0",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^organizations.*\",phase=\"Failed\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "0",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^users.*\",phase=\"Failed\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "0",
			}
			sum = append(sum, data)
		}
	}
	params = "?query=kube_pod_status_phase{pod=~\"^vtracking.*\",phase=\"Failed\"}"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := PodData{
				Pod:   result.Metric.Pod,
				Value: "0",
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetKpiVtrack(c *gin.Context) []KpiData {
	fmt.Println("check Get Kpi API")
	var response PodReponse
	var sum []KpiData
	var value, evalue string
	var ok, eok bool
	url := config.PrometheusUrl
	params := "?query=sum(users_api_request_error_count)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		evalue, eok = response.Data.Result[0].Value[1].(string)
		if !eok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			evalue = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		eok = false
		evalue = "0"
	}

	params = "?query=sum(users_api_request_count)"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		value, ok = response.Data.Result[0].Value[1].(string)
		if !ok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "1"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "1"
	}
	if ok && eok {
		fmt.Println("giá trị value " + value + " " + evalue)
		data := KpiData{
			Pod:   "user_services",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}
	if !ok || !eok {
		fmt.Println("khong co gia tri evalue ")
		data := KpiData{
			Pod:   "user_services",
			Req:   "1",
			Error: "0",
		}
		sum = append(sum, data)
	}

	params = "?query=sum(vtdevices_api_request_error_count)"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		evalue, eok = response.Data.Result[0].Value[1].(string)
		if !eok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			evalue = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		eok = false
		evalue = "0"
	}

	params = "?query=sum(vtdevices_api_request_count)"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		value, ok = response.Data.Result[0].Value[1].(string)
		if !ok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "1"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "1"
	}
	if ok && eok {
		fmt.Println("giá trị value " + value + " " + evalue)
		data := KpiData{
			Pod:   "devices_services",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}
	if !ok || !eok {
		fmt.Println("khong co gia tri evalue ")
		data := KpiData{
			Pod:   "devices_services",
			Req:   "1",
			Error: "0",
		}
		sum = append(sum, data)
	}
	params = "?query=sum(organizations_api_request_error_count)"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		evalue, eok = response.Data.Result[0].Value[1].(string)
		if !eok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			evalue = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		eok = false
		evalue = "0"
	}

	params = "?query=sum(organizations_api_request_count)"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		value, ok = response.Data.Result[0].Value[1].(string)
		if !ok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "1"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "1"
	}
	if ok && eok {
		fmt.Println("giá trị value " + value + " " + evalue)
		data := KpiData{
			Pod:   "organizations_services",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}
	if !ok || !eok {
		fmt.Println("khong co gia tri evalue ")
		data := KpiData{
			Pod:   "organizations_services",
			Req:   "1",
			Error: "0",
		}
		sum = append(sum, data)
	}

	params = "?query=sum(attributes_api_request_error_count)"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		evalue, eok = response.Data.Result[0].Value[1].(string)
		if !eok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			evalue = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		evalue = "0"
	}
	params = "?query=sum(attributes_api_request_count)"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		value, ok = response.Data.Result[0].Value[1].(string)
		if !ok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "0"
	}
	if ok && eok {
		data := KpiData{
			Pod:   "attribute_services",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}
	if !ok || !eok {
		fmt.Println("khong co gia tri evalue ")
		data := KpiData{
			Pod:   "attribute_services",
			Req:   "1",
			Error: "0",
		}
		sum = append(sum, data)
	}
	params = "?query=sum(vtracking_report_api_request_error_count{method=~\"VTrackingReportOverviewV4\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	fmt.Println("check 1")
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	fmt.Println("check 2")
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	fmt.Println("check 3")
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		evalue, eok = response.Data.Result[0].Value[1].(string)
		if !eok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			evalue = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		eok = false
		evalue = "0"
	}
	fmt.Println("check 4")
	params = "?query=sum(vtracking_report_api_request_count{method=~\"VTrackingReportOverviewV4\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	fmt.Println("check 5")
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	fmt.Println("check 6")
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	fmt.Println("check 8")
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		value, ok = response.Data.Result[0].Value[1].(string)
		if !ok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "1"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "1"
	}
	if ok && eok {
		fmt.Println("giá trị value " + value + " " + evalue)
		data := KpiData{
			Pod:   "vtracking_report_overview_api",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}
	if !ok || !eok {
		fmt.Println("khong co gia tri evalue ")
		data := KpiData{
			Pod:   "vtracking_report_overview_api",
			Req:   "1",
			Error: "0",
		}
		sum = append(sum, data)
	}
	fmt.Println("check 9")
	params = "?query=sum(vtracking_api_request_error_count{method=~\"get_images\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	fmt.Println("check 10")
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	fmt.Println("check 11")
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	fmt.Println("check 12")
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		evalue, eok = response.Data.Result[0].Value[1].(string)
		if !eok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			evalue = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		evalue = "0"
		eok = false
	}
	fmt.Println("check 13")
	params = "?query=sum(vtracking_api_request_count{method=~\"get_images\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	fmt.Println("check 14")
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	fmt.Println("check 15")
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		value, ok = response.Data.Result[0].Value[1].(string)
		if !ok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "1"
			ok = false
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "1"
	}
	if ok && eok {
		fmt.Println("giá trị value " + value + " " + evalue)
		data := KpiData{
			Pod:   "get_images_api",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}
	if !ok || !eok {
		fmt.Println("khong co gia tri value ")
		data := KpiData{
			Pod:   "get_images_api",
			Req:   "1",
			Error: "0",
		}
		sum = append(sum, data)
	}
	//fmt.Println("check 16")
	//params = "?query=sum(users_api_request_error_count{method=~\"user_vtracking_login\"})"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	c.Error(err)
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	c.Error(err)
	//}
	//
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	c.Error(err)
	//}
	//if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
	//	// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
	//	evalue, eok = response.Data.Result[0].Value[1].(string)
	//	if !eok {
	//		// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
	//		fmt.Println("Chưa có giá trị value")
	//		evalue = "0"
	//	}
	//} else {
	//	// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
	//	fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
	//	evalue = "0"
	//	eok = false
	//}
	//fmt.Println("check 17")
	//params = "?query=sum(users_api_request_count{method=~\"user_vtracking_login\"})"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	c.Error(err)
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	c.Error(err)
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	c.Error(err)
	//}
	//if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
	//	// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
	//	value, ok = response.Data.Result[0].Value[1].(string)
	//	if !ok {
	//		// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
	//		fmt.Println("Chưa có giá trị value")
	//		value = "0"
	//	}
	//} else {
	//	// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
	//	fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
	//	value = "0"
	//	ok = false
	//}
	//if ok && eok {
	//	fmt.Println("giá trị value " + value + " " + evalue)
	//	data := KpiData{
	//		Pod:   "user_login_api",
	//		Req:   value,
	//		Error: evalue,
	//	}
	//	sum = append(sum, data)
	//}
	//if !ok || !eok {
	//	fmt.Println("khong co gia tri value ")
	//	data := KpiData{
	//		Pod:   "user_login_api",
	//		Req:   "1",
	//		Error: "0",
	//	}
	//	sum = append(sum, data)
	//}
	fmt.Println("check 19")
	params = "?query=sum(attributes_api_request_error_count{method=~\"vtracking_get_attribute_time_series_paging\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		evalue, eok = response.Data.Result[0].Value[1].(string)
		if !eok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			evalue = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		evalue = "0"
		eok = false
	}
	fmt.Println("check 20")
	params = "?query=sum(attributes_api_request_count{method=~\"vtracking_get_attribute_time_series_paging\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		value, ok = response.Data.Result[0].Value[1].(string)
		if !ok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "1"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "1"
		ok = false
	}
	if ok && eok {
		fmt.Println("giá trị value " + value + " " + evalue)
		data := KpiData{
			Pod:   "vtracking_get_attribute_time_series_paging_api",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}
	if !ok || !eok {
		fmt.Println("giá trị value " + value + " " + evalue)
		data := KpiData{
			Pod:   "vtracking_get_attribute_time_series_paging_api",
			Req:   "1",
			Error: "0",
		}
		sum = append(sum, data)
	}
	//fmt.Println("check 21")
	//params = "?query=(sum(kube_pod_created-kube_pod_start_time)/sum(kube_pod_created-time()))"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	c.Error(err)
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	c.Error(err)
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	c.Error(err)
	//}
	//if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
	//	// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
	//	fvalue, fok := response.Data.Result[0].Value[1].(float64)
	//	roundedValue := math.Round(fvalue)
	//	value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
	//	fmt.Println("giá trị value " + value)
	//	data := KpiData{
	//		Pod:   "system_downtime",
	//		Req:   "1",
	//		Error: value,
	//	}
	//	if !fok {
	//		// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
	//		fmt.Println("Chưa có giá trị value")
	//		value = "0"
	//		data = KpiData{
	//			Pod:   "system_downtime",
	//			Req:   "1",
	//			Error: value,
	//		}
	//	}
	//	sum = append(sum, data)
	//} else {
	//	// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
	//	fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
	//	value = "0"
	//	//data := KpiData{
	//	//	Pod:   "system_downtime",
	//	//	Req:   "1",
	//	//	Error: value,
	//	//}
	//	//sum = append(sum, data)
	//}
	//
	//fmt.Println("check 22")
	//params = "?query=(sum(users_api_request_latency_microseconds_sum{method=\"user_vtracking_login\"})*1000)/sum(kube_pod_start_time{pod=~\"^users.*\"})"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	c.Error(err)
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	c.Error(err)
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	c.Error(err)
	//}
	//if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
	//	// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
	//	fvalue, fok := response.Data.Result[0].Value[1].(float64)
	//	roundedValue := math.Round(fvalue)
	//	value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
	//	fmt.Println("giá trị value " + value)
	//	data := KpiData{
	//		Pod:   "user_login_latency",
	//		Req:   "1",
	//		Error: value,
	//	}
	//
	//	if !fok {
	//		// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
	//		fmt.Println("Chưa có giá trị value")
	//		value = "0"
	//		data = KpiData{
	//			Pod:   "user_login_latency",
	//			Req:   "1",
	//			Error: value,
	//		}
	//	}
	//	sum = append(sum, data)
	//} else {
	//	// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
	//	fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
	//	value = "0"
	//	fmt.Println("giá trị value " + value)
	//	//data := KpiData{
	//	//	Pod:   "user_login_latency",
	//	//	Req:   "1",
	//	//	Error: value,
	//	//}
	//	//sum = append(sum, data)
	//}
	//fmt.Println("check 23")
	//params = "?query=(sum(vtracking_report_api_request_latency_microseconds_sum{method=\"VTrackingReportOverviewV3\"})*1000)/sum(kube_pod_start_time{pod=~\"^vtracking.*\"})"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	c.Error(err)
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	c.Error(err)
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	c.Error(err)
	//}
	//if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
	//	// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
	//	fvalue, fok := response.Data.Result[0].Value[1].(float64)
	//	roundedValue := math.Round(fvalue)
	//	value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
	//	fmt.Println("giá trị value " + value)
	//	data := KpiData{
	//		Pod:   "report_overview_latency",
	//		Req:   "1",
	//		Error: value,
	//	}
	//	if !fok {
	//		// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
	//		fmt.Println("Chưa có giá trị value")
	//		value = "0"
	//		data = KpiData{
	//			Pod:   "report_overview_latency",
	//			Req:   "1",
	//			Error: value,
	//		}
	//	}
	//	sum = append(sum, data)
	//} else {
	//	// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
	//	fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
	//	value = "0"
	//	//data := KpiData{
	//	//	Pod:   "user_login_latency",
	//	//	Req:   "1",
	//	//	Error: value,
	//	//}
	//	//sum = append(sum, data)
	//}
	//
	//params = "?query=(sum(vtracking_api_request_latency_microseconds_sum{method=\"get_images\"})*1000)/sum(kube_pod_start_time{pod=~\"^vtracking.*\"})"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	c.Error(err)
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	c.Error(err)
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	c.Error(err)
	//}
	//if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
	//	// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
	//	fvalue, fok := response.Data.Result[0].Value[1].(float64)
	//	roundedValue := math.Round(fvalue)
	//	value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
	//	fmt.Println("giá trị value " + value)
	//	data := KpiData{
	//		Pod:   "get_images_latency",
	//		Req:   "1",
	//		Error: value,
	//	}
	//	if !fok {
	//		// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
	//		fmt.Println("Chưa có giá trị value")
	//		value = "0"
	//		data = KpiData{
	//			Pod:   "get_images_latency",
	//			Req:   "1",
	//			Error: value,
	//		}
	//	}
	//	sum = append(sum, data)
	//} else {
	//	// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
	//	fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
	//	value = "0"
	//	//data := KpiData{
	//	//	Pod:   "user_login_latency",
	//	//	Req:   "1",
	//	//	Error: value,
	//	//}
	//	//sum = append(sum, data)
	//}
	//
	//params = "?query=(sum(attributes_api_request_latency_microseconds_sum{method=\"vtracking_get_attribute_time_series_paging\"})*1000)/sum(kube_pod_start_time{pod=~\"^attributes.*\"})"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	c.Error(err)
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	c.Error(err)
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	c.Error(err)
	//}
	//if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
	//	// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
	//	fvalue, fok := response.Data.Result[0].Value[1].(float64)
	//	roundedValue := math.Round(fvalue)
	//	value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
	//	fmt.Println("giá trị value " + value)
	//	data := KpiData{
	//		Pod:   "get_attribute_time_series_latency",
	//		Req:   "1",
	//		Error: value,
	//	}
	//	if !fok {
	//		// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
	//		fmt.Println("Chưa có giá trị value")
	//		value = "0"
	//		data = KpiData{
	//			Pod:   "get_attribute_time_series_latency",
	//			Req:   "1",
	//			Error: value,
	//		}
	//	}
	//	sum = append(sum, data)
	//} else {
	//	// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
	//	fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
	//	value = "0"
	//	//data := KpiData{
	//	//	Pod:   "user_login_latency",
	//	//	Req:   "1",
	//	//	Error: value,
	//	//}
	//	//sum = append(sum, data)
	//}
	return sum

}

func GetRequestKPIReport(c *gin.Context) KpiData {
	var response PodReponse
	var data KpiData
	var value, evalue string
	var ok, eok bool
	url := config.PrometheusUrl
	params := "?query=sum(vtracking_report_api_request_error_count{method=~\"VTrackingReportOverviewV4\"})"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	fmt.Println("check 1")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	fmt.Println("check 2")
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	fmt.Println("check 3")
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		evalue, eok = response.Data.Result[0].Value[1].(string)
		if !eok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			evalue = "0"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		eok = false
		evalue = "0"
	}
	fmt.Println("check 4")
	params = "?query=sum(vtracking_report_api_request_count{method=~\"VTrackingReportOverviewV4\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	fmt.Println("check 5")
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	fmt.Println("check 6")
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	fmt.Println("check 8")
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		value, ok = response.Data.Result[0].Value[1].(string)
		if !ok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "1"
		}
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "1"
	}
	if ok && eok {
		fmt.Println("giá trị value " + value + " " + evalue)
		data = KpiData{
			Pod:   "vtracking_report_overview_api",
			Req:   value,
			Error: evalue,
		}
	}
	if !ok || !eok {
		fmt.Println("khong co gia tri evalue ")
		data = KpiData{
			Pod:   "vtracking_report_overview_api",
			Req:   "1",
			Error: "0",
		}
	}

	return data
}

func GetDownTimeVtrack(c *gin.Context) []KpiData {
	var response PodReponse
	var sum []KpiData
	var value string
	url := config.PrometheusUrl

	params := "?query=(sum(kube_pod_created-kube_pod_start_time)/sum(kube_pod_created-time()))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		fvalue, fok := response.Data.Result[0].Value[1].(float64)
		roundedValue := math.Round(fvalue)
		value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
		fmt.Println("giá trị value " + value)
		data := KpiData{
			Pod:   "system_downtime",
			Req:   "1",
			Error: value,
		}
		if !fok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "0"
			data = KpiData{
				Pod:   "system_downtime",
				Req:   "1",
				Error: "0",
			}
		}
		sum = append(sum, data)
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "0"
		data := KpiData{
			Pod:   "system_downtime",
			Req:   "1",
			Error: "0",
		}
		sum = append(sum, data)
	}

	fmt.Println("check 22")
	return sum
}

func GetKpiLatencyVtrack(c *gin.Context) []KpiData {
	fmt.Println("check Get Kpi API")
	var response PodReponse
	var sum []KpiData
	var value string
	url := config.PrometheusUrl

	params := "?query=(sum(kube_pod_created-kube_pod_start_time)/sum(kube_pod_created-time()))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		fvalue, fok := response.Data.Result[0].Value[1].(float64)
		roundedValue := math.Round(fvalue)
		value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
		fmt.Println("giá trị value " + value)
		data := KpiData{
			Pod:   "system_downtime",
			Req:   "1",
			Error: value,
		}
		if !fok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "0"
			data = KpiData{
				Pod:   "system_downtime",
				Req:   "1",
				Error: value,
			}
		}
		sum = append(sum, data)
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "0"
		//data := KpiData{
		//	Pod:   "system_downtime",
		//	Req:   "1",
		//	Error: value,
		//}
		//sum = append(sum, data)
	}

	fmt.Println("check 22")
	params = "?query=(sum(users_api_request_latency_microseconds_sum{method=\"user_vtracking_login\"})*1000)/sum(kube_pod_start_time{pod=~\"^users.*\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		fvalue, fok := response.Data.Result[0].Value[1].(float64)
		roundedValue := math.Round(fvalue)
		value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
		fmt.Println("giá trị value " + value)
		data := KpiData{
			Pod:   "user_login_latency",
			Req:   "1",
			Error: value,
		}

		if !fok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "0"
			data = KpiData{
				Pod:   "user_login_latency",
				Req:   "1",
				Error: value,
			}
		}
		sum = append(sum, data)
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "0"
		fmt.Println("giá trị value " + value)
		//data := KpiData{
		//	Pod:   "user_login_latency",
		//	Req:   "1",
		//	Error: value,
		//}
		//sum = append(sum, data)
	}
	//fmt.Println("check 23")
	//params = "?query=(sum(vtracking_report_api_request_latency_microseconds_sum{method=\"VTrackingReportOverviewV3\"})*1000)/sum(kube_pod_start_time{pod=~\"^vtracking.*\"})"
	//resp, err = http.Get(url + params)
	//if err != nil {
	//	log.Printf("error in services %s", err)
	//	c.Error(err)
	//}
	//
	//body, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("error reading response body: %s", err)
	//	c.Error(err)
	//}
	//err = json.Unmarshal(body, &response)
	//if err != nil {
	//	log.Printf("error unmarshaling JSON: %s", err)
	//	c.Error(err)
	//}
	//if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
	//	// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
	//	fvalue, fok := response.Data.Result[0].Value[1].(float64)
	//	roundedValue := math.Round(fvalue)
	//	value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
	//	fmt.Println("giá trị value " + value)
	//	data := KpiData{
	//		Pod:   "report_overview_latency",
	//		Req:   "1",
	//		Error: value,
	//	}
	//	if !fok {
	//		// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
	//		fmt.Println("Chưa có giá trị value")
	//		value = "0"
	//		data = KpiData{
	//			Pod:   "report_overview_latency",
	//			Req:   "1",
	//			Error: value,
	//		}
	//	}
	//	sum = append(sum, data)
	//} else {
	//	// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
	//	fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
	//	value = "0"
	//	//data := KpiData{
	//	//	Pod:   "user_login_latency",
	//	//	Req:   "1",
	//	//	Error: value,
	//	//}
	//	//sum = append(sum, data)
	//}

	params = "?query=(sum(vtracking_api_request_latency_microseconds_sum{method=\"get_images\"})*1000)/sum(kube_pod_start_time{pod=~\"^vtracking.*\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		fvalue, fok := response.Data.Result[0].Value[1].(float64)
		roundedValue := math.Round(fvalue)
		value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
		fmt.Println("giá trị value " + value)
		data := KpiData{
			Pod:   "get_images_latency",
			Req:   "1",
			Error: value,
		}
		if !fok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "0"
			data = KpiData{
				Pod:   "get_images_latency",
				Req:   "1",
				Error: value,
			}
		}
		sum = append(sum, data)
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "0"
		//data := KpiData{
		//	Pod:   "user_login_latency",
		//	Req:   "1",
		//	Error: value,
		//}
		//sum = append(sum, data)
	}

	params = "?query=(sum(attributes_api_request_latency_microseconds_sum{method=\"vtracking_get_attribute_time_series_paging\"})*1000)/sum(kube_pod_start_time{pod=~\"^attributes.*\"})"
	resp, err = http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	if len(response.Data.Result) > 0 && len(response.Data.Result[0].Value) > 1 {
		// Kiểm tra xem mảng Value trả về có ít nhất 2 phần tử hay không
		fvalue, fok := response.Data.Result[0].Value[1].(float64)
		roundedValue := math.Round(fvalue)
		value = strconv.FormatFloat(roundedValue, 'f', -1, 64)
		fmt.Println("giá trị value " + value)
		data := KpiData{
			Pod:   "get_attribute_time_series_latency",
			Req:   "1",
			Error: value,
		}
		if !fok {
			// Giá trị không tồn tại hoặc không thể chuyển đổi thành kiểu string
			fmt.Println("Chưa có giá trị value")
			value = "0"
			data = KpiData{
				Pod:   "get_attribute_time_series_latency",
				Req:   "1",
				Error: value,
			}
		}
		sum = append(sum, data)
	} else {
		// Không có giá trị nào trong mảng Value hoặc không đủ phần tử để truy cập
		fmt.Println("Không có giá trị trong mảng Value hoặc không đủ phần tử")
		value = "0"
		//data := KpiData{
		//	Pod:   "user_login_latency",
		//	Req:   "1",
		//	Error: value,
		//}
		//sum = append(sum, data)
	}
	return sum

}

func GetPrometheus(c *gin.Context, query string) interface{} {
	var response map[string]interface{}
	baseURL, err := url.Parse(config.PrometheusUrl)
	if err != nil {
		log.Printf("error parsing prometheus url: %s", err)
		c.Error(err)
		return ""
	}

	params := baseURL.Query()
	params.Set("query", query)
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
		return ""
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
		return ""
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := string(body)
		if len(msg) > 300 {
			msg = msg[:300]
		}
		err = fmt.Errorf("prometheus returned status=%d body=%q", resp.StatusCode, msg)
		log.Printf("error in services %s", err)
		c.Error(err)
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		msg := string(body)
		if len(msg) > 300 {
			msg = msg[:300]
		}
		log.Printf("error unmarshaling JSON: %s body=%q", err, msg)
		c.Error(err)
		return ""
	}

	return response

}

func GetMqttRequest(c *gin.Context) interface{} {
	var response KafkaReponse
	url := config.PrometheusUrl
	params := "?query=sum(rate(mqtt_publish_received[1m]))"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	return response.Data.Result[0].Value[1]
}

func GetMqttClientConnected(c *gin.Context) interface{} {
	var response KafkaReponse
	url := config.PrometheusUrl
	var sum []UrlData
	params := "?query=(socket_open-socket_close)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := UrlData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func RedisMemUsed(c *gin.Context) interface{} {
	var response KafkaReponse
	url := config.PrometheusUrl
	var sum []UrlData
	params := "?query=redis_memory_used_bytes"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := UrlData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func RedisCmdPS(c *gin.Context) interface{} {
	var response CmdReponse
	url := config.PrometheusUrl
	var sum []CmdData
	params := "?query=sum(rate(redis_commands_total{}[1m]))by(cmd)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := CmdData{
				Cmd:   result.Metric.Cmd,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetInodeUsage(c *gin.Context) interface{} {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=node_filesystem_avail_bytes{}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := SysData{
				Url:    result.Metric.Instance,
				Device: result.Metric.Device,
				Value:  value,
			}
			sum = append(sum, data)
		}
	}

	return sum
}

func GetReadOnlyFile(c *gin.Context) interface{} {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=node_filesystem_readonly{}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}

	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := SysData{
				Url:    result.Metric.Instance,
				Device: result.Metric.Device,
				Value:  value,
			}
			sum = append(sum, data)
		}
	}

	return sum
}

func GetCpuUsage(c *gin.Context) interface{} {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=100-ssCpuIdle{}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && result.Metric.Job == "snmp kemp_loadmaster" {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetMemTotal(c *gin.Context) []SysData {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=memTotalReal"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && result.Metric.Job == "snmp kemp_loadmaster" {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetMemAvail(c *gin.Context) []SysData {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=memAvailReal"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && result.Metric.Job == "snmp kemp_loadmaster" {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetMemCache(c *gin.Context) interface{} {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=memCached"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && result.Metric.Job == "snmp kemp_loadmaster" {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetMemBuffer(c *gin.Context) interface{} {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=memBuffer"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && result.Metric.Job == "snmp kemp_loadmaster" {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetIoReceive(c *gin.Context) []SysData {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=ssIOReceive"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetIoSent(c *gin.Context) []SysData {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=ssIOSent"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func DiskTotal(c *gin.Context) []SysData {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=hrStorageSize{hrStorageDescr=\"/\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func DiskUsed(c *gin.Context) []SysData {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=hrStorageUsed{hrStorageDescr=\"/\"}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func MaxfileOpen(c *gin.Context) []SysData {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=node_filefd_maximum{}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func NTPSynchronous(c *gin.Context) []SysData {
	var response Systemreponse
	var sum []SysData
	url := config.PrometheusUrl
	params := "?query=node_timex_sync_status{}"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.Error(err)
	}
	defer resp.Body.Close() // Đảm bảo body được đóng sau khi sử dụng.

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %s", err)
		c.Error(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
		c.Error(err)
	}
	for _, result := range response.Data.Result {
		if value, ok := result.Value[1].(string); ok && value == "1" {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: "Syn",
			}
			sum = append(sum, data)
		}
		if value, ok := result.Value[1].(string); ok && value == "0" {
			data := SysData{
				Url:   result.Metric.Instance,
				Value: "NotSyn",
			}
			sum = append(sum, data)
		}
	}
	return sum
}
