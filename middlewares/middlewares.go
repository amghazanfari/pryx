package middlewares

import (
	"github.com/gin-gonic/gin"
)

func CheckSuperAdmin(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		gotToken := c.Request.Header.Get("Authorization")

		if token == gotToken {
			c.Next()
		} else {
			c.JSON(403, gin.H{
				"error": "permission denied",
			})
			c.Abort()
			return
		}
	}
}
