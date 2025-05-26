package router

import (
	"os"
	"profile/internal/handlers"
	"profile/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRouter(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://pollforge.ru"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	profile_port := os.Getenv("PORT")
	protected := r.Group("/api/profile")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/", handlers.GetProfile)
		protected.GET("/user/:id", handlers.GetDifUserProfile)
		protected.GET("/forms", handlers.GetForms)
		protected.PUT("/name", handlers.UpdateProfileName)
		protected.PUT("/bio", handlers.UpdateProfileBio)

		protected.DELETE("/", handlers.DeleteProfile)
		protected.DELETE("/forms/:id", handlers.DeleteForm)
		protected.POST("/avatar", handlers.UploadAvatar)
	}

	if err := r.Run(profile_port); err != nil {
		panic(err)
	}
}
