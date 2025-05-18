package router

import (
	"forms/internal/handlers"
	"forms/internal/middleware"
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
	forms_port := os.Getenv("PORT")
	protected := r.Group("/api/forms")

	r.GET("/forms/link/:link", handlers.GetFormByLink).Use(middleware.JWTAuth())

	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/", handlers.CreateForm)
		protected.GET("/:id", handlers.GetForm)
		protected.PUT("/:id", handlers.UpdateForm)
		protected.DELETE("/:id", handlers.DeleteForm)

		protected.POST("/:id/questions", handlers.CreateQuestion)
		protected.PUT("/:id/questions/:question_id", handlers.UpdateQuestion)
		protected.DELETE("/:id/questions/:question_id", handlers.DeleteQuestion)

		protected.POST("/:id/questions/:question_id/answers", handlers.CreateAnswer)
		protected.PUT("/:id/questions/:question_id/answers/:answer_id", handlers.UpdateAnswer)
		protected.DELETE("/:id/questions/:question_id/answers/:answer_id", handlers.DeleteAnswer)
	}
	if err := r.Run(forms_port); err != nil {
		panic(err)
	}
}
