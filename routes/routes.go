package routes

import (
	"Monitor_Platform/controllers"
	"Monitor_Platform/middlewares"
	"Monitor_Platform/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	route.Use(cors.New(config))

	route.Use(func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if err.IsType(gin.ErrorTypeBind) {
					// Trả về lỗi 500
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
					return
				}
			}
			// Lấy lỗi đầu tiên trong danh sách lỗi
			err := c.Errors[0].Err
			// Trả về lỗi
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})

	route.Use(middlewares.ErrorHandler())
	route.GET("ems/snmp/oid/get", middlewares.ValidateOIDMiddleware, controllers.GetOIDController)
	route.GET("ems/snmp/server/get", middlewares.ValidateOidServerMiddleware, controllers.GetOidWithServerController)
	route.GET("ems/snmp/memory/total", controllers.GetAllMemController)
	route.GET("ems/snmp/memory/available", controllers.GetAvailMemController)
	route.GET("ems/snmp/memory/used", controllers.GetUsedMemController)
	route.GET("ems/snmp/memory/cache", controllers.GetMemCacheController)
	route.GET("ems/snmp/memory/buffer", controllers.GetMemBufferController)
	route.GET("ems/snmp/memory/all", controllers.GetMemController)
	route.GET("ems/snmp/memory/swap", controllers.GetMemSwapController)
	route.GET("ems/snmp/disk/all", controllers.GetDiskController)
	route.GET("ems/snmp/cpu/usage", controllers.GetCpuUsageController)
	route.GET("ems/snmp/io/data", controllers.GetIODataController)
	route.GET("ems/users/tps", controllers.GetUsersTPSController)
	route.GET("ems/organizations/tps", controllers.GetOrganizationsTPSController)
	route.GET("ems/adapters/tps", controllers.GetAdapterTPSController)
	route.GET("ems/roles/tps", controllers.GetRoleTPSController)
	route.GET("ems/devices/tps", controllers.GetDeviceTPSController)
	route.GET("ems/video/stream", controllers.GetDataUsageStreamingController)
	route.GET("ems/url/state", controllers.GetUrlStateController)
	route.GET("ems/http/status/code", controllers.GetHttpStatusController)
	route.GET("ems/pod/state", controllers.GetPodStateController)
	route.GET("ems/kafka/queue", controllers.GetKafkaQueueController)
	route.GET("ems/kafka/partition", controllers.GetKafkaPartitionController)
	route.GET("ems/pod/state/all", controllers.GetPodupController)
	route.GET("ems/pod/pending/all", controllers.GetPodpendController)
	route.GET("ems/kpi/vtrack", controllers.GetKpiVtrackController)
	route.GET("ems/redis/cmd/rate", controllers.GetCmdRedisController)
	route.GET("ems/redis/mem/used", controllers.GetRedisMemController)
	route.GET("ems/system/inode", controllers.GetInodeController)
	route.GET("ems/system/file/rxonly", controllers.GetRxOnlyFile)
	route.GET("ems/mqtt/client/connect", controllers.GetMqttClientController)
	route.GET("ems/mqtt/message/receive", controllers.GetRateMqttMessageController)
	route.GET("ems/prometheus", middlewares.ValidateQueryMiddleware, controllers.GetPrometheusController)

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
