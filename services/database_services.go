package services

import (
	"Monitor_Platform/config"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func GetOpsMongo(c *gin.Context) []MongoData {
	var response Systemreponse
	var sum []MongoData
	url := config.PrometheusUrl
	params := "?query=rate(mongodb_op_counters_total{}[30m])"
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
			data := MongoData{
				Url:   result.Metric.Instance,
				Type:  result.Metric.Type,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetDocumentMongo(c *gin.Context) []MongoData {
	var response Systemreponse
	var sum []MongoData
	url := config.PrometheusUrl
	params := "?query=rate(mongodb_mongod_metrics_document_total{}[30m])"
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
			data := MongoData{
				Url:   result.Metric.Instance,
				Type:  result.Metric.Type,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetHealthyMem(c *gin.Context) []MongoData {
	var response Systemreponse
	var sum []MongoData
	url := config.PrometheusUrl
	params := "?query=mongodb_mongod_replset_number_of_members{}"
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
			data := MongoData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetConnectNum(c *gin.Context) []MongoData {
	var response Systemreponse
	var sum []MongoData
	url := config.PrometheusUrl
	params := "?query=mongodb_connections{state=\"current\"}"
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
			data := MongoData{
				Url:   result.Metric.Instance,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetMemMongo(c *gin.Context) []MongoData {
	var response Systemreponse
	var sum []MongoData
	url := config.PrometheusUrl
	params := "?query=mongodb_memory{type=~\"resident|virtual\"}"
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
			data := MongoData{
				Url:   result.Metric.Instance,
				Type:  result.Metric.Type,
				Value: value,
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func GetNetworkMongo(c *gin.Context) []MongoData {
	var response Systemreponse
	var sum []MongoData
	url := config.PrometheusUrl
	params := "?query=rate(mongodb_network_metrics_num_requests_total{}[30m])"
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
			data := MongoData{
				Url:   result.Metric.Instance,
				Value: value, //values is bytes
			}
			sum = append(sum, data)
		}
	}
	return sum
}
