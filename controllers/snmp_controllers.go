package controllers

import (
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

func GetAllCpuUsageController(c *gin.Context) {
	values := services.GetAllCpuUsage(c)
	c.JSON(200, values)
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
	values := services.GetDataSNMP(server, oid)
	c.JSON(200, values)
}
