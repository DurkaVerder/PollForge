package router

import (
	"auth/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {
	auth := r.Group("/api")
	{
		auth.POST("/register", handlers.UserRegistration)
		auth.POST("/logging", handlers.UserLogging)
	}
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
