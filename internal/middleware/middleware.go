package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// Logger middleware
func Logger() gin.HandlerFunc {
	return gin.Logger()
}

// RequireAuth middleware to protect authenticated routes
func RequireAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get session/token from cookie or header
		session := c.GetHeader("Authorization")
		if session == "" {
			// Try to get from cookie
			cookie, err := c.Cookie("session")
			if err != nil || cookie == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}
			session = cookie
		}

		// For now, we'll implement a simple session check
		// In a real app, you'd validate the JWT token or session ID
		if session == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// Set user context (you'd decode the actual user from token/session)
		// c.Set("user_id", userID)

		c.Next()
	})
}
