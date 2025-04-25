package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware(c gin.Context) {
	// Check exists userID in the jwt

	c.Next()
}
