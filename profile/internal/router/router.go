package router

import (
	"profile/internal/handlers"
	"profile/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.POST("/forms", handlers.GetForms)
	}
	if err := r.Run(":8084"); err != nil {
		panic(err)
	}
}
