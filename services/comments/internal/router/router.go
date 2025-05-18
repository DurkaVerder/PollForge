package router

import (
	"comments/internal/handlers"
	"comments/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRouter(r *gin.Engine) {

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	auth_port := os.Getenv("PORT")
	protected := r.Group("/api/comments")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/forms/:form_id/comments", handlers.GetComments)
		protected.POST("/forms/:form_id/comments", handlers.CreateComment)
		protected.PUT("/forms/:form_id/comments/:id", handlers.UpdateComment)
		protected.DELETE("/forms/:form_id/comments/:id", handlers.DeleteComment)
	}
	if err := r.Run(auth_port); err != nil {
		panic(err)
	}
}
