package controllers

import (
	"Monitor_Platform/services"
	"Monitor_Platform/validations"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

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
	c.JSON(200, data)
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
			Rate: num / req,
		}
		reponse = append(reponse, cal)
	}

	c.JSON(200, reponse)

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
