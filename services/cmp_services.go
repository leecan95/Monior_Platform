package services

import (
	"Monitor_Platform/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func BccsTPS(c *gin.Context) interface{} {
	url := config.PrometheusUrl
	params := "?query = round(sum(irate(telcojob_bccs_request_count{}[30s])) by (gateway), 0.001)"
	resp, err := http.Get(url + params)
	if err != nil {
		log.Printf("error in services %s", err)
		c.JSON(400, err)
	}

	return resp
}
