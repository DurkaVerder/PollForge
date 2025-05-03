package router

import (
	"os"
	"profile/internal/handlers"
	"profile/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {
	profile_port := os.Getenv("PORT")
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.GET("/forms", handlers.GetForms)
		protected.DELETE("/forms/:id", handlers.DeleteForm)

	}
	
	if err := r.Run(profile_port); err != nil {
		panic(err)
	}
}
