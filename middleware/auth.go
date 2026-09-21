package middleware

import (
	"fmt"
	"net/http"

	auth "github.com/JerryJeager/ohara-be/internal/http"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := auth.ValidateToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     "Bad request",
				"message":    "Authentication failed",
				"statusCode": http.StatusUnauthorized,
			})
			fmt.Println(err)
			c.Abort()
			return
		}
		c.Set("user_id", id)

		c.Next()
	}
}


func RefreshAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := auth.GetTokenFromRequest(c)
		id, err := auth.ValidateToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     "Bad request",
				"message":    "Authentication failed",
				"statusCode": http.StatusUnauthorized,
			})
			fmt.Println(err)
			c.Abort()
			return
		}
		c.Set("user_id", id)
		c.Set("refresh_token", refreshToken)

		c.Next()
	}
}
