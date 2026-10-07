package middleware

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {

	allowedOrigins := map[string]bool{
		"http://localhost:3000":       true,
		"http://localhost:5173":       true,
		"https://ohara-v0.vercel.app": true,
	}

	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		origin := c.Request.Header.Get("Origin")

		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
