package router

import (
	"auth/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {
	auth_port := os.Getenv("PORT")
	auth := r.Group("/api")
	{
		auth.POST("/register", handlers.UserRegistration)
		auth.POST("/logging", handlers.UserLogging)
	}
	if err := r.Run(auth_port); err != nil {
		panic(err)
	}
}
