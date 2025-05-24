package router

import (
	"auth/internal/handlers"
	"os"
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
	auth_port := os.Getenv("PORT")
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.UserRegistration)
		auth.POST("/logging", handlers.UserLogging)
		auth.POST("/password_resets", handlers.PasswordResetRequest)
		auth.POST("/password_resets/confirm", handlers.PasswordResetConfirmHandler)
	}
	if err := r.Run(auth_port); err != nil {
		panic(err)
	}
}
