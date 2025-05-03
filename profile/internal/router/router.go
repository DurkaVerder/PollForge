package router

import (
	"os"
	"profile/internal/handlers"
	"profile/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRouter(r *gin.Engine) {

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	profile_port := os.Getenv("PORT")
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.GET("/forms", handlers.GetForms)
	}

	if err := r.Run(profile_port); err != nil {
		panic(err)
	}
}
