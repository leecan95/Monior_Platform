package controllers

import (
	"Monitor_Platform/services"
	"Monitor_Platform/validations"
	"github.com/gin-gonic/gin"
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

func GetPodStateController(c *gin.Context) {
	values := services.GetPodState(c)
	if values != "" {
		c.JSON(200, values)
	}
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
