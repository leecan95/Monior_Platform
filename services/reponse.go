package services

import (
	"github.com/gin-gonic/gin"
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
