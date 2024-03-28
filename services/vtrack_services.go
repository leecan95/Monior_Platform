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
	var response PodReponse
	var sum []KpiData
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
	evalue, eok := response.Data.Result[0].Value[1].(string)

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

	value, ok := response.Data.Result[0].Value[1].(string)
	if ok && eok {
		data := KpiData{
			Pod:   "user",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}
	//
	//params = "?query=sum(vtdevices_api_request_error_count)"
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
	//
	//evalue, eok = response.Data.Result[0].Value[1].(string)
	//
	//params = "?query=sum(vtdevices_api_request_count)"
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
	//value, ok = response.Data.Result[0].Value[1].(string)
	//if ok && eok {
	//	data := KpiData{
	//		Pod:   "device",
	//		Req:   value,
	//		Error: evalue,
	//	}
	//	sum = append(sum, data)
	//}
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

	evalue, eok = response.Data.Result[0].Value[1].(string)

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
	value, ok = response.Data.Result[0].Value[1].(string)
	if ok && eok {
		data := KpiData{
			Pod:   "organization",
			Req:   value,
			Error: evalue,
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

	evalue, eok = response.Data.Result[0].Value[1].(string)

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
	value, ok = response.Data.Result[0].Value[1].(string)
	if ok && eok {
		data := KpiData{
			Pod:   "attribute",
			Req:   value,
			Error: evalue,
		}
		sum = append(sum, data)
	}

	return sum
}

func GetPrometheus(c *gin.Context, query string) interface{} {
	var response map[string]interface{}
	url := config.PrometheusUrl
	params := "?query=" + query
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
		return ""
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("error unmarshaling JSON: %s", err)
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
