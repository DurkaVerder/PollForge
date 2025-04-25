package middleware

import (
	"net/http"
	"stream_line/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {

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
