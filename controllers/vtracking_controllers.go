package controllers

import (
	"Monitor_Platform/services"
	"Monitor_Platform/validations"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetUsersTPSController(c *gin.Context) {
	values := services.GetUsersTPS(c)
	c.JSON(200, values)
}

func GetOrganizationsTPSController(c *gin.Context) {
	values := services.GetOrganizationTPS(c)
	c.JSON(200, values)
}

func GetAdapterTPSController(c *gin.Context) {
	values := services.GetAdapterTPS(c)
	c.JSON(200, values)
}

func GetRoleTPSController(c *gin.Context) {
	values := services.GetRoleTPS(c)
	c.JSON(200, values)
}

func GetDeviceTPSController(c *gin.Context) {
	values := services.GetDeviceTPS(c)
	c.JSON(200, values)
}

func GetDataUsageStreamingController(c *gin.Context) {
	values := services.GetDataUsageStreaming(c)
	c.JSON(200, values)
}

func GetUrlStateController(c *gin.Context) {
	values := services.GetUrlState(c)
	c.JSON(200, values)
}

func GetPodStateController(c *gin.Context) {
	values := services.GetPodState(c)
	c.JSON(200, values)
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
	c.JSON(200, values)
}
