package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
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

/*** PostgreSql ***/

func Maxconnections(c *gin.Context) []MongoData {
	var response Systemreponse
	var sum []MongoData
	url := config.PrometheusUrl
	params := "?query=pg_settings_max_connections"
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

func DatabaseSize(c *gin.Context) []MongoData {
	var response Systemreponse
	var sum []MongoData
	url := config.PrometheusUrl
	params := "?query=sum(pg_database_size_bytes{})"
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

func PgTPSCommit(c *gin.Context) []PgData {
	var response Pgreponse
	var sum []PgData

	url := config.PrometheusUrl
	params := "?query=irate(pg_stat_database_xact_commit{datname=~\"organizations|users|devices|vtracking1|attributes|vtracking1_0\"}[5m])"
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
			data := PgData{
				Url:       result.Metric.Instance,
				Namespace: result.Metric.Namespace,
				Value:     value, //values is bytes
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func PgTPSRollback(c *gin.Context) []PgData {
	var response Pgreponse
	var sum []PgData

	url := config.PrometheusUrl
	params := "?query=irate(pg_stat_database_xact_rollback{datname=~\"organizations|users|devices|vtracking1|attributes|vtracking1_0\"}[5m])"
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
			data := PgData{
				Url:       result.Metric.Instance,
				Namespace: result.Metric.Namespace,
				Value:     value, //values is bytes
			}
			sum = append(sum, data)
		}
	}
	return sum
}

func PgLatencyKpi(c *gin.Context) []config.LatencyKpi {
	var data config.LatencyKpi
	var reponse []config.LatencyKpi

	data, _ = model.GetOverallLatency(c)
	reponse = append(reponse, data)
	data, _ = model.GetReportLatency(c)
	reponse = append(reponse, data)
	data, _ = model.GetLoginLatency(c)
	reponse = append(reponse, data)
	data, _ = model.GetImageLatency(c)
	reponse = append(reponse, data)
	return reponse
}

func PgLatencyTotalKpi(c *gin.Context) config.LatencyKpi {
	var data config.LatencyKpi
	data, _ = model.GetOverallLatency(c)
	return data
}

func PgLatencyLoginKpi(c *gin.Context) config.LatencyKpi {
	var data config.LatencyKpi
	data, _ = model.GetLoginLatency(c)
	return data
}

func PgLatencyReportKpi(c *gin.Context) config.LatencyKpi {
	var data config.LatencyKpi
	data, _ = model.GetReportLatency(c)
	return data
}

func PgLatencyGetImagesKpi(c *gin.Context) config.LatencyKpi {
	var data config.LatencyKpi
	data, _ = model.GetImageLatency(c)
	return data
}

func PgLatencyTrackingKpi(c *gin.Context) config.LatencyKpi {
	var data config.LatencyKpi
	data, _ = model.GetTrackingLatency(c)
	return data
}

func PgRequestKpi(c *gin.Context) []config.SuccessKpi {
	var data config.SuccessKpi
	var reponse []config.SuccessKpi

	data, _ = model.GetRequestTotal(c)
	reponse = append(reponse, data)
	data, _ = model.GetRequestLoginTotal(c)
	reponse = append(reponse, data)
	data, _ = model.GetRequestGetImageTotal(c)
	reponse = append(reponse, data)
	data, _ = model.GetRequestReportTotal(c)
	reponse = append(reponse, data)

	return reponse
}

func PgRequestTotalKpi(c *gin.Context) config.SuccessKpi {
	var data config.SuccessKpi
	data, _ = model.GetRequestTotal(c)
	return data
}

func PgRequestLoginKpi(c *gin.Context) config.SuccessKpi {
	var data config.SuccessKpi
	data, _ = model.GetRequestLoginTotal(c)
	return data
}

func PgRequestGetImageKpi(c *gin.Context) config.SuccessKpi {
	var data config.SuccessKpi
	data, _ = model.GetRequestGetImageTotal(c)
	return data
}
func PgRequestReportKpi(c *gin.Context) config.SuccessKpi {
	var data config.SuccessKpi
	data, _ = model.GetRequestReportTotal(c)
	return data
}

func PgRequestTrackingKpi(c *gin.Context) config.SuccessKpi {
	var data config.SuccessKpi
	data, _ = model.GetRequestTrackingTotal(c)
	return data
}
