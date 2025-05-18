package router

import (
	"auth/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRouter(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth_port := os.Getenv("PORT")
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.UserRegistration)
		auth.POST("/logging", handlers.UserLogging)
	}
	if err := r.Run(auth_port); err != nil {
		panic(err)
	}
}
