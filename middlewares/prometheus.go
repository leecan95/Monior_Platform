package middlewares

import "github.com/gin-gonic/gin"

func ValidateQueryMiddleware(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(400, gin.H{
			"message": "query parameters is required",
		})
		c.Abort()
		return
	}
	c.Next()
}
