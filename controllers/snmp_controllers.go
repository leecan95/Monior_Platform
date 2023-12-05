package controllers

import (
	"Monitor_Platform/config"
	"Monitor_Platform/services"
	"Monitor_Platform/validations"
	"github.com/gin-gonic/gin"
	"strings"
)

var previousIdle = make([]int64, 6)
var previousUser = make([]int64, 6)

func GetOIDController(c *gin.Context) {
	var request validations.GetOidSnmp
	if requestErr := c.ShouldBind(&request); requestErr != nil {
		errRes := strings.Split(requestErr.Error(), ": ")
		services.ValidationError(c, "These fields are required!", errRes)
		return
	}
	// Lấy danh sách các OID từ query parameters
	oids := c.QueryArray("oid")
	// Thực hiện lấy giá trị tương ứng với từng OID
	values := services.GetValueFromOID(c, oids)

	c.JSON(200, values)

}

func GetAllMemController(c *gin.Context) {
	values := services.GetAllMemInfo(c)
	data := []MemSnmpAll{
		{Server: config.Server244, Value: values[0]},
		{Server: config.Server245, Value: values[1]},
		{Server: config.Server246, Value: values[2]},
		{Server: config.Server250, Value: values[3]},
		{Server: config.Server251, Value: values[4]},
		{Server: config.Server252, Value: values[5]},
	}
	c.JSON(200, data)
}

func GetAvailMemController(c *gin.Context) {
	values := services.GetAvailMemInfo(c)
	data := []MemSnmpAll{
		{Server: config.Server244, Value: values[0]},
		{Server: config.Server245, Value: values[1]},
		{Server: config.Server246, Value: values[2]},
		{Server: config.Server250, Value: values[3]},
		{Server: config.Server251, Value: values[4]},
		{Server: config.Server252, Value: values[5]},
	}
	c.JSON(200, data)
}

func GetUsedMemController(c *gin.Context) {
	total := services.GetAllMemInfo(c)
	avail := services.GetAvailMemInfo(c)
	cache := services.GetMemCacheInfo(c)
	buffer := services.GetMemBufferInfo(c)
	data := []MemSnmpAll{
		{Server: config.Server244, Value: total[0] - (avail[0] + cache[0] + buffer[0])},
		{Server: config.Server245, Value: total[1] - (avail[1] + cache[1] + buffer[1])},
		{Server: config.Server246, Value: total[2] - (avail[2] + cache[2] + buffer[2])},
		{Server: config.Server250, Value: total[3] - (avail[3] + cache[3] + buffer[3])},
		{Server: config.Server251, Value: total[4] - (avail[4] + cache[4] + buffer[4])},
		{Server: config.Server252, Value: total[5] - (avail[5] + cache[5] + buffer[5])},
	}
	c.JSON(200, data)
}

func GetMemController(c *gin.Context) {
	total := services.GetAllMemInfo(c)
	avail := services.GetAvailMemInfo(c)
	cache := services.GetMemCacheInfo(c)
	buffer := services.GetMemBufferInfo(c)
	data := []MemAll{
		{Server: config.Server244, Total: total[0], Available: avail[0] + cache[0] + buffer[0]},
		{Server: config.Server245, Total: total[1], Available: avail[1] + cache[1] + buffer[1]},
		{Server: config.Server246, Total: total[2], Available: avail[2] + cache[2] + buffer[2]},
		{Server: config.Server250, Total: total[3], Available: avail[3] + cache[3] + buffer[3]},
		{Server: config.Server251, Total: total[4], Available: avail[4] + cache[4] + buffer[4]},
		{Server: config.Server252, Total: total[5], Available: avail[5] + cache[5] + buffer[5]},
	}
	c.JSON(200, data)
}

func GetMemCacheController(c *gin.Context) {
	cache := services.GetMemCacheInfo(c)
	data := []MemSnmpAll{
		{Server: config.Server244, Value: cache[0]},
		{Server: config.Server245, Value: cache[1]},
		{Server: config.Server246, Value: cache[2]},
		{Server: config.Server250, Value: cache[3]},
		{Server: config.Server251, Value: cache[4]},
		{Server: config.Server252, Value: cache[5]},
	}
	c.JSON(200, data)
}

func GetMemBufferController(c *gin.Context) {
	buffer := services.GetMemBufferInfo(c)
	data := []MemSnmpAll{
		{Server: config.Server244, Value: buffer[0]},
		{Server: config.Server245, Value: buffer[1]},
		{Server: config.Server246, Value: buffer[2]},
		{Server: config.Server250, Value: buffer[3]},
		{Server: config.Server251, Value: buffer[4]},
		{Server: config.Server252, Value: buffer[5]},
	}
	c.JSON(200, data)
}

func GetCpuUsageVerController(c *gin.Context) {
	user := services.GetCpuPercentUser(c)
	sys := services.GetCpuPercentSystem(c)
	data := []CpuSnmpAll{
		{Server: config.Server244, Value: float64(user[0] + sys[0])},
		{Server: config.Server245, Value: float64(user[1] + sys[1])},
		{Server: config.Server246, Value: float64(user[2] + sys[2])},
		{Server: config.Server250, Value: float64(user[3] + sys[3])},
		{Server: config.Server251, Value: float64(user[4] + sys[4])},
		{Server: config.Server252, Value: float64(user[5] + sys[5])},
	}
	c.JSON(200, data)
}

func GetCpuUsageController(c *gin.Context) {
	user := services.GetCpuRawUser(c)
	idle := services.GetCpuRawIdle(c)
	var data []CpuSnmpAll
	//if previousUser[0] == 0 || previousUser[1] == 0 || previousUser[2] == 0 || previousUser[3] == 0 || previousUser[4] == 0 || previousUser[5] == 0 {
	data = []CpuSnmpAll{
		{Server: config.Server244, Value: float64((user[0] * 100) / (user[0] + idle[0]))},
		{Server: config.Server245, Value: float64((user[1] * 100) / (user[1] + idle[1]))},
		{Server: config.Server246, Value: float64((user[2] * 100) / (user[2] + idle[2]))},
		{Server: config.Server250, Value: float64((user[3] * 100) / (user[3] + idle[3]))},
		{Server: config.Server251, Value: float64((user[4] * 100) / (user[4] + idle[4]))},
		{Server: config.Server252, Value: float64((user[5] * 100) / (user[5] + idle[5]))},
	}
	//} else {
	//	data = []CpuSnmpAll{
	//		{Server: config.Server244, Value: float64(((user[0] - previousUser[0]) * 100) / ((user[0] - previousUser[0]) + (idle[0] - previousIdle[0])))},
	//		{Server: config.Server245, Value: float64(((user[1] - previousUser[1]) * 100) / ((user[1] - previousUser[1]) + (idle[1] - previousIdle[1])))},
	//		{Server: config.Server246, Value: float64(((user[2] - previousUser[2]) * 100) / ((user[2] - previousUser[2]) + (idle[2] - previousIdle[2])))},
	//		{Server: config.Server250, Value: float64(((user[3] - previousUser[3]) * 100) / ((user[3] - previousUser[3]) + (idle[3] - previousIdle[3])))},
	//		{Server: config.Server251, Value: float64(((user[4] - previousUser[4]) * 100) / ((user[4] - previousUser[4]) + (idle[4] - previousIdle[4])))},
	//		{Server: config.Server252, Value: float64(((user[5] - previousUser[5]) * 100) / ((user[5] - previousUser[5]) + (idle[5] - previousIdle[5])))},
	//	}
	//}

	previousUser[0] = user[0]
	previousUser[1] = user[1]
	previousUser[2] = user[2]
	previousUser[3] = user[3]
	previousUser[4] = user[4]
	previousUser[5] = user[5]
	previousIdle[0] = idle[0]
	previousIdle[1] = idle[1]
	previousIdle[2] = idle[2]
	previousIdle[3] = idle[3]
	previousIdle[4] = idle[4]
	previousIdle[5] = idle[5]
	c.JSON(200, data)
}

func GetDiskController(c *gin.Context) {
	total := services.GetDiskTotal(c)
	avail := services.GetDiskAvail(c)
	data := []MemAll{
		{Server: config.Server244, Total: total[0], Available: avail[0]},
		{Server: config.Server245, Total: total[1], Available: avail[1]},
		{Server: config.Server246, Total: total[2], Available: avail[2]},
		{Server: config.Server250, Total: total[3], Available: avail[3]},
		{Server: config.Server251, Total: total[4], Available: avail[4]},
		{Server: config.Server252, Total: total[5], Available: avail[5]},
	}
	c.JSON(200, data)
}

func GetIODataController(c *gin.Context) {
	receive := services.GetIoReceiveData(c)
	sent := services.GetIoSentData(c)
	data := []IoData{
		{Server: config.Server244, Receive: receive[0], Sent: sent[0]},
		{Server: config.Server245, Receive: receive[1], Sent: sent[1]},
		{Server: config.Server246, Receive: receive[2], Sent: sent[2]},
		{Server: config.Server250, Receive: receive[3], Sent: sent[3]},
		{Server: config.Server251, Receive: receive[4], Sent: sent[4]},
		{Server: config.Server252, Receive: receive[5], Sent: sent[5]},
	}
	c.JSON(200, data)
}
func GetMemSwapController(c *gin.Context) {
	avail := services.GetMemSwapAvail(c)
	total := services.GetMemSwapTotal(c)
	data := []MemAll{
		{Server: config.Server244, Total: total[0], Available: avail[0]},
		{Server: config.Server245, Total: total[1], Available: avail[1]},
		{Server: config.Server246, Total: total[2], Available: avail[2]},
		{Server: config.Server250, Total: total[3], Available: avail[3]},
		{Server: config.Server251, Total: total[4], Available: avail[4]},
		{Server: config.Server252, Total: total[5], Available: avail[5]},
	}
	c.JSON(200, data)
}

func GetOidWithServerController(c *gin.Context) {
	var request validations.GetEachOidWithServer
	if requestErr := c.ShouldBind(&request); requestErr != nil {
		errRes := strings.Split(requestErr.Error(), ": ")
		services.ValidationError(c, "These fields are required!", errRes)
		return
	}
	oid := c.Query("oid")
	server := c.Query("server")
	values := services.GetDataSNMP(c, server, oid)
	c.JSON(200, values)
}
