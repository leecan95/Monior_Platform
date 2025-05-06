package controllers

import (
	"Monitor_Platform/config"
	"Monitor_Platform/services"
	"github.com/gin-gonic/gin"
)
func GetCMPKpiLatencyDBController(c *gin.Context) {
	var data []config.LatencyKpi
	data = services.PgCMPLatencyKpi(c)
	c.JSON(200, data)
}