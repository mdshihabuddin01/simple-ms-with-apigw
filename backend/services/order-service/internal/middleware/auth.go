package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
		// 	c.Abort()
		// 	return
		// }

		// if strings.HasPrefix(authHeader, "Bearer ") {
		// 	// tokenString = authHeader[7:]
		// }

		userID := c.GetHeader("X-User-ID")
		// userEmail := c.GetHeader("X-User-Email")
		// userActive := c.GetHeader("X-User-Active")

		// if userID == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context headers"})
		// 	c.Abort()
		// 	return
		// }

		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context headers"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
