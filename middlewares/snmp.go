package middlewares

import "github.com/gin-gonic/gin"

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
