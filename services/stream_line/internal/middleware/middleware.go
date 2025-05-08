package middleware

import (
	"log"
	"net/http"
	"stream_line/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwt := c.GetHeader("Authorization")
		if jwt == "" {
			c.Set("userID", "-1") // No userID in the jwt, set to -1
			c.Next()
		}

		userID, err := service.GetParamFromJWT(jwt, "id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
		}

		c.Set("userID", userID)

		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := log.New(c.Writer, "[Middleware]", log.LstdFlags|log.Lshortfile)
		logger.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		logger.Printf("Response: %d", c.Writer.Status())
	}
}
