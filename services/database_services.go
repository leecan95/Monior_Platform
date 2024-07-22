package services

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

func ExportKpiToExcel() {
	var data config.SuccessKpi
	var d config.LatencyKpi
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	sheetName := "Sheet1"
	f.SetSheetName(f.GetSheetName(0), sheetName)
	headers := []string{"Services", "Rate"}
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
		style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
		f.SetCellStyle(sheetName, cell, cell, style)
	}
	data, _ = model.ExportRequestTotal()
	f.SetCellValue("Sheet1", "A2", "total services")
	f.SetCellValue("Sheet1", "B2", data.Value)
	data, _ = model.ExportRequestLoginTotal()
	f.SetCellValue("Sheet1", "A3", "login services")
	f.SetCellValue("Sheet1", "B3", data.Value)
	data, _ = model.ExportRequestGetImageTotal()
	f.SetCellValue("Sheet1", "A4", "get images services")
	f.SetCellValue("Sheet1", "B4", data.Value)
	data, _ = model.ExportRequestReportTotal()
	f.SetCellValue("Sheet1", "A4", "report services")
	f.SetCellValue("Sheet1", "B4", data.Value)
	data, _ = model.ExportRequestTrackingTotal()
	f.SetCellValue("Sheet1", "A5", "tracking services")
	f.SetCellValue("Sheet1", "B5", data.Value)
	d, _ = model.ExportOverallLatency()
	f.SetCellValue("Sheet1", "A5", "Latency KPI Overall")
	style, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellStyle(sheetName, "A5", "A5", style)
	f.SetCellValue("Sheet1", "A6", "Total request")
	f.SetCellValue("Sheet1", "B6", d.Total)
	f.SetCellValue("Sheet1", "A7", "Number of requests with a response time less than 5s")
	f.SetCellValue("Sheet1", "B7", d.Count)
	f.SetCellValue("Sheet1", "A8", "Percentile Request 5s")
	f.SetCellValue("Sheet1", "B8", d.Percentile)
	d, _ = model.ExportLoginLatency()
	f.SetCellValue("Sheet1", "A9", "Latency KPI Login")
	style, _ = f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellStyle(sheetName, "A9", "A9", style)
	f.SetCellValue("Sheet1", "A10", "Total request")
	f.SetCellValue("Sheet1", "B10", d.Total)
	f.SetCellValue("Sheet1", "A11", "Number of requests with a response time less than 5s")
	f.SetCellValue("Sheet1", "B11", d.Count)
	f.SetCellValue("Sheet1", "A12", "Percentile Request 5s")
	f.SetCellValue("Sheet1", "B12", d.Percentile)

	d, _ = model.ExportImageLatency()
	f.SetCellValue("Sheet1", "A13", "Latency KPI Get Image")
	style, _ = f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellStyle(sheetName, "A13", "A13", style)
	f.SetCellValue("Sheet1", "A14", "Total request")
	f.SetCellValue("Sheet1", "B14", d.Total)
	f.SetCellValue("Sheet1", "A15", "Number of requests with a response time less than 5s")
	f.SetCellValue("Sheet1", "B15", d.Count)
	f.SetCellValue("Sheet1", "A16", "Percentile Request 5s")
	f.SetCellValue("Sheet1", "B16", d.Percentile)
	d, _ = model.ExportTrackingLatency()
	f.SetCellValue("Sheet1", "A17", "Latency KPI Tracking")
	style, _ = f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellStyle(sheetName, "A17", "A17", style)
	f.SetCellValue("Sheet1", "A18", "Total request")
	f.SetCellValue("Sheet1", "B18", d.Total)
	f.SetCellValue("Sheet1", "A19", "Number of requests with a response time less than 5s")
	f.SetCellValue("Sheet1", "B19", d.Count)
	f.SetCellValue("Sheet1", "A20", "Percentile Request 5s")
	f.SetCellValue("Sheet1", "B20", d.Percentile)
	d, _ = model.ExportReportLatency()
	f.SetCellValue("Sheet1", "A21", "Latency KPI Report")
	style, _ = f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	f.SetCellStyle(sheetName, "A21", "A21", style)
	f.SetCellValue("Sheet1", "A22", "Total request")
	f.SetCellValue("Sheet1", "B22", d.Total)
	f.SetCellValue("Sheet1", "A23", "Number of requests with a response time less than 5s")
	f.SetCellValue("Sheet1", "B24", d.Count)
	f.SetCellValue("Sheet1", "A24", "Percentile Request 5s")
	f.SetCellValue("Sheet1", "B24", d.Percentile)

	// Lấy thời gian hiện tại
	now := time.Now()

	// Định dạng thời gian thành chuỗi theo định dạng ddmmyyyy
	formattedTime := now.Format("02012006")

	// Tạo tên file với thời gian
	fileName := fmt.Sprintf("KPI_Vtracking_%s.xlsx", formattedTime)

	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("File saved successfully at: %s\n", fileName) // Log the file path
	}
}
