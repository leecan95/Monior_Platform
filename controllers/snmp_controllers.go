package controllers

import (
	"Monitor_Platform/config"
	"Monitor_Platform/services"
	"Monitor_Platform/validations"
	"github.com/gin-gonic/gin"
	"strings"
)

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
	data := []MemSnmpAll{
		{Server: config.Server244, Value: total[0] - avail[0]},
		{Server: config.Server245, Value: total[1] - avail[1]},
		{Server: config.Server246, Value: total[2] - avail[2]},
		{Server: config.Server250, Value: total[3] - avail[3]},
		{Server: config.Server251, Value: total[4] - avail[4]},
		{Server: config.Server252, Value: total[5] - avail[5]},
	}
	c.JSON(200, data)
}

func GetMemController(c *gin.Context) {
	total := services.GetAllMemInfo(c)
	avail := services.GetAvailMemInfo(c)
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

func GetCpuUsageController(c *gin.Context) {
	user := services.GetCpuRawUser(c)
	sys := services.GetCpuRawSystem(c)
	nice := services.GetCpuRawNice(c)
	idle := services.GetCpuRawIdle(c)
	irq := services.GetCpuRawInterrupt(c)
	soft := services.GetCpuRawSoftIrq(c)
	rawK := services.GetCpuRawKernel(c)
	rawW := services.GetCpuRawWait(c)
	data := []CpuSnmpAll{
		{Server: config.Server244, Value: float64(((user[0] + sys[0] + nice[0] + irq[0] + soft[0] + rawW[0] + rawK[0]) * 100) / (user[0] + sys[0] + nice[0] + irq[0] + soft[0] + rawW[0] + rawK[0] + idle[0]))},
		{Server: config.Server245, Value: float64(((user[1] + sys[1] + nice[1] + irq[1] + soft[1] + rawW[1] + rawK[1]) * 100) / (user[1] + sys[1] + nice[1] + irq[1] + soft[1] + rawW[1] + rawK[1] + idle[1]))},
		{Server: config.Server246, Value: float64(((user[2] + sys[2] + nice[2] + irq[2] + soft[2] + rawW[2] + rawK[2]) * 100) / (user[2] + sys[2] + nice[2] + irq[2] + soft[2] + rawW[2] + rawK[2] + idle[2]))},
		{Server: config.Server250, Value: float64(((user[3] + sys[3] + nice[3] + irq[3] + soft[3] + rawW[3] + rawK[3]) * 100) / (user[3] + sys[3] + nice[3] + irq[3] + soft[3] + rawW[3] + rawK[3] + idle[3]))},
		{Server: config.Server251, Value: float64(((user[4] + sys[4] + nice[4] + irq[4] + soft[4] + rawW[4] + rawK[4]) * 100) / (user[4] + sys[4] + nice[4] + irq[4] + soft[4] + rawW[4] + rawK[4] + idle[4]))},
		{Server: config.Server252, Value: float64(((user[5] + sys[5] + nice[5] + irq[5] + soft[5] + rawW[5] + rawK[5]) * 100) / (user[5] + sys[5] + nice[5] + irq[5] + soft[5] + rawW[5] + rawK[5] + idle[5]))},
	}
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
