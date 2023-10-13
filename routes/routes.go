package routes

import (
	"Monitor_Platform/controllers"
	"Monitor_Platform/middlewares"
	"Monitor_Platform/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	route.Use(cors.New(config))

	route.GET("snmp/oid/get", middlewares.ValidateOIDMiddleware, controllers.GetOIDController)
	route.GET("snmp/server/get", middlewares.ValidateOidServerMiddleware, controllers.GetOidWithServerController)
	route.GET("users/tps", controllers.GetUsersTPSController)
	// Endpoint to get metrics
	// Middleware để ghi log vào metric
	route.Use(func(c *gin.Context) {
		services.SetCPUTemperature(66.8) // Giả định giá trị nhiệt độ CPU
		c.Next()
	})
	route.Use(func(c *gin.Context) {
		services.MemToTalAvail()
		c.Next()
	})
	route.Use(func(c *gin.Context) {
		services.IncHttpRequest(c)
		c.Next()
	})
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return route
}
