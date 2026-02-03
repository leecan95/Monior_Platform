package middlewares

import "github.com/gin-gonic/gin"
import "net/http"
import "strconv"

const maxLimit = 1000

func ValidateTransactionQueryMiddleware(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > maxLimit {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid limit",
		})
		c.Abort()
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid offset",
		})
		c.Abort()
		return
	}

	action := c.Query("action")
	if len(action) > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid action",
		})
		c.Abort()
		return
	}

	// WHY: parse 1 lần, dùng nhiều nơi
	c.Set("limit", limit)
	c.Set("offset", offset)
	c.Set("action", action)

	c.Next()
}
