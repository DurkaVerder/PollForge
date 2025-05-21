package router

import (
	"admin/internal/handlers"
	"admin/internal/middleware"
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
	admin_port := os.Getenv("PORT")
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
	{
		admin.GET("/users", handlers.ListUsers)
		admin.PUT("/users/:id/ban", handlers.ToggleBanUser)
		admin.DELETE("/users/:id", handlers.DeleteUser)

		admin.DELETE("/forms/:id", handlers.DeleteForm)

		admin.DELETE("/comments/:id", handlers.DeleteComment)

	}
	if err := r.Run(admin_port); err != nil {
		panic(err)
	}
}
