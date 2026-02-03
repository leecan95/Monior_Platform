package controllers

import (
	"Monitor_Platform/config"
	"Monitor_Platform/model"
	"Monitor_Platform/services"
	"Monitor_Platform/validations"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

//NEW

type ApiLogRequest struct {
	Url        string `json:"url" binding:"required"`
	StatusCode int    `json:"status_code" binding:"required"`
	Latency    int64  `json:"latency" binding:"required"`
	Timestamp  int64  `json:"timestamp"`
}

func CollectApiLogController(c *gin.Context) {
	var req ApiLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logItem := model.ApiLog{
		Url:        req.Url,
		StatusCode: req.StatusCode,
		Latency:    req.Latency,
		Timestamp:  time.UnixMilli(req.Timestamp),
	}

	services.EnqueueApiLog(logItem)

	c.JSON(200, gin.H{"message": "accepted"})
}

// NEWWWW
func GetUsersTPSController(c *gin.Context) {
	values := services.GetUsersTPS(c)

	if values != "" {
		reponse := Response{
			Value: values,
		}
		c.JSON(200, reponse)
	}
}

func GetOrganizationsTPSController(c *gin.Context) {
	values := services.GetOrganizationTPS(c)
	if values != "" {
		reponse := Response{
			Value: values,
		}
		c.JSON(200, reponse)
	}
}

func GetAdapterTPSController(c *gin.Context) {
	values := services.GetAdapterTPS(c)
	if values != "" {
		reponse := Response{
			Value: values,
		}
		c.JSON(200, reponse)
	}
}

func GetRoleTPSController(c *gin.Context) {
	values := services.GetRoleTPS(c)
	if values != "" {
		reponse := Response{
			Value: values,
		}
		c.JSON(200, reponse)
	}
}

func GetDeviceTPSController(c *gin.Context) {
	values := services.GetDeviceTPS(c)
	if values != "" {
		reponse := Response{
			Value: values,
		}
		c.JSON(200, reponse)
	}
}

func GetDataUsageStreamingController(c *gin.Context) {
	values := services.GetDataUsageStreaming(c)
	if values != "" {
		c.JSON(200, values)
	}
}

func GetUrlStateController(c *gin.Context) {
	values := services.GetUrlState(c)
	if values != "" {
		c.JSON(200, values)
	}
}

func GetHttpStatusController(c *gin.Context) {
	values := services.GetHttpStatusCode(c)
	if values != "" {
		c.JSON(200, values)
	}
}

func GetPodStateController(c *gin.Context) {
	values := services.GetPodState(c)
	c.JSON(200, values)
}

func GetKafkaQueueController(c *gin.Context) {
	var byteOut, byteIn services.KafkaData
	var KafkaQueue []Kafka
	numIn := make([]float64, 0, len(byteIn.Result))
	numOut := make([]float64, 0, len(byteOut.Result))
	byteOut = services.GetKafkaByteout(c)
	byteIn = services.GetKafkaBytein(c)
	for _, result := range byteIn.Result {
		//num := int64(result.Value[1].(string))
		num, err := strconv.ParseFloat(result.Value[1].(string), 64)
		if err != nil {
			fmt.Println("Error converting value:", err)
			return
		}
		numIn = append(numIn, num)
	}
	for _, result := range byteOut.Result {
		num, err := strconv.ParseFloat(result.Value[1].(string), 64)
		if err != nil {
			fmt.Println("Error converting value:", err)
			return
		}
		numOut = append(numOut, num)
	}
	for i, _ := range byteIn.Result {
		data := Kafka{
			Server:  byteIn.Result[i].Metric.Instance,
			Bytein:  numIn[i],
			Byteout: numOut[i],
		}
		KafkaQueue = append(KafkaQueue, data)

	}
	c.JSON(200, KafkaQueue)
}

func GetKafkaPartitionController(c *gin.Context) {
	data := services.GetKafkaPartitionOnline(c)

	reponse := Response{
		Value: data,
	}
	c.JSON(200, reponse)
}

func GetPodupController(c *gin.Context) {
	data := services.GetPodUp(c)
	c.JSON(200, data)
}

func GetPodpendController(c *gin.Context) {
	data := services.GetPodPend(c)
	c.JSON(200, data)
}

func GetKpiVtrackController(c *gin.Context) {
	var data []services.KpiData
	var reponse []KpiVtrack
	data = services.GetKpiVtrack(c)

	for _, result := range data {
		req, err := strconv.ParseFloat(result.Req, 64)
		num, err1 := strconv.ParseFloat(result.Error, 64)
		if err != nil {
			c.Error(err)
		}
		if err1 != nil {
			c.Error(err1)
		}
		cal := KpiVtrack{
			Pod:  result.Pod,
			Rate: (100 - (num/req)*100),
		}
		reponse = append(reponse, cal)
	}

	c.JSON(200, reponse)

}

func GetKpiLatencyVtrackController(c *gin.Context) {
	var data []services.KpiData
	var reponse []KpiVtrack
	data = services.GetKpiLatencyVtrack(c)

	for _, result := range data {
		req, err := strconv.ParseFloat(result.Req, 64)
		num, err1 := strconv.ParseFloat(result.Error, 64)
		if err != nil {
			c.Error(err)
		}
		if err1 != nil {
			c.Error(err1)
		}
		cal := KpiVtrack{
			Pod:  result.Pod,
			Rate: (num / req) * 100,
		}
		reponse = append(reponse, cal)
	}

	c.JSON(200, reponse)

}

func GetKpiDowntimeController(c *gin.Context) {
	var data []services.KpiData
	var reponse []KpiVtrack
	data = services.GetDownTimeVtrack(c)

	for _, result := range data {
		req, err := strconv.ParseFloat(result.Req, 64)
		num, err1 := strconv.ParseFloat(result.Error, 64)
		if err != nil {
			c.Error(err)
		}
		if err1 != nil {
			c.Error(err1)
		}
		cal := KpiVtrack{
			Pod:  result.Pod,
			Rate: (num / req) * 100,
		}
		reponse = append(reponse, cal)
	}

	c.JSON(200, reponse)
}

func GetKpiLatencyDBController(c *gin.Context) {
	var data []config.LatencyKpi
	data = services.PgLatencyKpi(c)
	c.JSON(200, data)
}

func GetKpiLatencyDBTotalController(c *gin.Context) {
	var data config.LatencyKpi
	data = services.PgLatencyTotalKpi(c)
	c.JSON(200, data)
}

func GetKpiLatencyDBLoginController(c *gin.Context) {
	var data config.LatencyKpi
	data = services.PgLatencyLoginKpi(c)
	c.JSON(200, data)
}

func GetKpiLatencyDBReportController(c *gin.Context) {
	var data config.LatencyKpi
	data = services.PgLatencyReportKpi(c)
	c.JSON(200, data)
}

func GetKpiLatencyDBGetImageController(c *gin.Context) {
	var data config.LatencyKpi
	data = services.PgLatencyGetImagesKpi(c)
	c.JSON(200, data)
}

func GetKpiLatencyDBTrackingController(c *gin.Context) {
	var data config.LatencyKpi
	data = services.PgLatencyTrackingKpi(c)
	c.JSON(200, data)
}

func GetKpiRequestDBController(c *gin.Context) {
	var data []config.SuccessKpi
	data = services.PgRequestKpi(c)
	c.JSON(200, data)
}

func GetKpiRequestDBTotalController(c *gin.Context) {
	var data config.SuccessKpi
	data = services.PgRequestTotalKpi(c)
	c.JSON(200, data)
}

func GetKpiRequestDBLoginController(c *gin.Context) {
	var data config.SuccessKpi
	data = services.PgRequestLoginKpi(c)
	c.JSON(200, data)
}

func GetKpiRequestDBGetImagesController(c *gin.Context) {
	var data config.SuccessKpi
	data = services.PgRequestGetImageKpi(c)
	c.JSON(200, data)
}
func ResolveActionURL(action string) (string, bool) {
	url, ok := config.ActionURLMap[action]
	return url, ok
}
func GetTransactionsController(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	action := c.Query("action")

	var url string
	if action != "" {
		u, ok := ResolveActionURL(action)
		if !ok {
			c.JSON(400, gin.H{"error": "invalid action"})
			return
		}
		url = u
	}

	data, err := services.GetTransactions(
		c,
		url,
		offset,
		limit,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, data)
}

func GetLatencyPercentileByURLController(c *gin.Context) {
	action := c.Query("action")
	var url string
	if action != "" {
		u, ok := ResolveActionURL(action)
		if !ok {
			c.JSON(400, gin.H{"error": "invalid action"})
			return
		}
		url = u
	}

	if url == "" {
		c.JSON(400, gin.H{
			"error": "url query parameter is required",
		})
		return
	}

	result, err := services.GetLatencyPercentileByURLService(c, url)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, result)
}

func GetKpiRequestDBReportController(c *gin.Context) {
	var data config.SuccessKpi
	//data = services.PgRequestReportKpi(c)
	var result services.KpiData
	result = services.GetRequestKPIReport(c)
	req, err := strconv.ParseFloat(result.Req, 64)
	num, err1 := strconv.ParseFloat(result.Error, 64)
	success := req - num
	if err != nil {
		c.Error(err)
	}
	if err1 != nil {
		c.Error(err1)
	}
	data = config.SuccessKpi{
		Api:   "Report",
		Value: (success / req) * 100,
	}
	c.JSON(200, data)
}

func GetKpiRequestDBTrackingController(c *gin.Context) {
	var data config.SuccessKpi
	data = services.PgRequestTrackingKpi(c)
	c.JSON(200, data)
}

func GetCmdRedisController(c *gin.Context) {
	data := services.RedisCmdPS(c)
	c.JSON(200, data)
}

func GetRedisMemController(c *gin.Context) {
	data := services.RedisMemUsed(c)
	c.JSON(200, data)
}

func GetRateMqttMessageController(c *gin.Context) {
	data := services.GetMqttRequest(c)
	reponse := Response{
		Value: data,
	}
	c.JSON(200, reponse)
}

func GetMqttClientController(c *gin.Context) {
	data := services.GetMqttClientConnected(c)
	c.JSON(200, data)
}
func GetInodeController(c *gin.Context) {
	data := services.GetInodeUsage(c)
	c.JSON(200, data)
}

func GetRxOnlyFile(c *gin.Context) {
	data := services.GetReadOnlyFile(c)
	c.JSON(200, data)
}

//func GetDiskProController(c *gin.Context) {
//	var total, used []services.SysData
//	var mem []MemAll
//	total = services.DiskTotal(c)
//	used = services.DiskUsed(c)
//	var tnum, unum int
//	var err error
//	for _, tdata := range total {
//		for _, udata := range used {
//			if tdata.Url == udata.Url {
//				tnum, err = strconv.Atoi(tdata.Value)
//				if err != nil {
//					c.Error(err)
//				}
//				unum, err = strconv.Atoi(udata.Value)
//				if err != nil {
//					c.Error(err)
//				}
//				data := MemAll{
//					Server:    tdata.Url,
//					Total:     tdata.Value,
//					Available: strconv.Itoa(tnum - unum),
//				}
//				mem = append(mem, data)
//			}
//		}
//	}
//	c.JSON(200, mem)
//}

//func GetMemProController(c *gin.Context) {
//	var total, avail []services.SysData
//	var mem []MemAll
//	total = services.GetMemTotal(c)
//	avail = services.GetMemAvail(c)
//	for _, tdata := range total {
//		for _, udata := range avail {
//			if tdata.Url == udata.Url {
//				data := MemAll{
//					Server:    tdata.Url,
//					Total:     tdata.Value,
//					Available: udata.Value,
//				}
//				mem = append(mem, data)
//			}
//		}
//	}
//	c.JSON(200, mem)
//}

func GetIoProController(c *gin.Context) {
	var iosend, ioreceive []services.SysData
	var iodata []IoData
	iosend = services.GetIoSent(c)
	ioreceive = services.GetIoReceive(c)
	for _, sdata := range iosend {
		for _, rdata := range ioreceive {
			if sdata.Url == rdata.Url {
				data := IoData{
					Server:  sdata.Url,
					Sent:    sdata.Value,
					Receive: rdata.Value,
				}
				iodata = append(iodata, data)
			}
		}
	}
	c.JSON(200, iodata)
}

func GetMemCacheProController(c *gin.Context) {
	data := services.GetMemCache(c)
	c.JSON(200, data)
}

func GetMemBufferProController(c *gin.Context) {
	data := services.GetMemBuffer(c)
	c.JSON(200, data)
}
func GetPrometheusController(c *gin.Context) {
	var request = validations.GetPrometheusVtrack{}
	if requestErr := c.ShouldBind(&request); requestErr != nil {
		errRes := strings.Split(requestErr.Error(), ": ")
		services.ValidationError(c, "These fields are required!", errRes)
		return
	}
	query := c.Query("query")

	values := services.GetPrometheus(c, query)

	if values != "" {
		c.JSON(200, values)
	}
}
