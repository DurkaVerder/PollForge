package router

import (
	"auth/internal/handlers"
	"auth/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {
	auth := r.Group("/api")
	{
		auth.POST("/register", handlers.UserRegistration)
		auth.POST("/logging", handlers.UserLogging)
	}
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.POST("/forms", handlers.GetForms)
	}
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
