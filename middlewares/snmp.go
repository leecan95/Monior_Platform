package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ValidateOIDMiddleware(c *gin.Context) {
	oid := c.Query("oid")
	if oid == "" {
		c.JSON(400, gin.H{
			"message": "oid parameters is required",
		})
		c.Abort()
		return
	}
	c.Next()
}

func ValidateOidServerMiddleware(c *gin.Context) {
	oid := c.Query("oid")
	if oid == "" {
		c.JSON(400, gin.H{
			"message": "oid parameters is required",
		})
		c.Abort()
		return
	}
	server := c.Query("server")
	if server == "" {
		c.JSON(400, gin.H{
			"message": "server parameters is required",
		})
		c.Abort()
		return
	}
	c.Next()
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in ErrorHandler", r)
				c.JSON(500, gin.H{
					"error": "Internal Server Error",
				})
			}
		}()
		c.Next()
	}
}
