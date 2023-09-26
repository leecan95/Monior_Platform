package services

import (
	"Monitor_Platform/snmp"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func SuccessResponse(c *gin.Context, status int, message interface{}, data interface{}) {
	c.AbortWithStatusJSON(status, gin.H{"message": message, "data": data})
}
func Success(c *gin.Context, message interface{}, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

func ErrorResponse(c *gin.Context, status int, message interface{}, error interface{}) {
	c.AbortWithStatusJSON(status, gin.H{"message": message, "errors": error})
}
func ValidationError(c *gin.Context, message interface{}, error interface{}) {
	ErrorResponse(c, http.StatusUnprocessableEntity, message, error)
}
func GetValueFromOID(c *gin.Context, oid []string) interface{} {
	result, err2 := snmp.Manager.Get(oid) // Lấy dữ liệu từ SNMP agent
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}
	return result
}
