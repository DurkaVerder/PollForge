package router

import (
	"forms/internal/handlers"
	"forms/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRouter(r *gin.Engine) {

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	forms_port := os.Getenv("PORT")
	protected := r.Group("/api")

	r.GET("/forms/link/:link", handlers.GetFormByLink).Use(middleware.JWTAuth())

	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/forms", handlers.CreateForm)
		protected.GET("/forms/:id", handlers.GetForm)
		protected.PUT("/forms/:id", handlers.UpdateForm)
		protected.DELETE("/forms/:id", handlers.DeleteForm)

		protected.POST("/forms/:id/questions", handlers.CreateQuestion)
		protected.PUT("/forms/:id/questions/:question_id", handlers.UpdateQuestion)
		protected.DELETE("/forms/:id/questions/:question_id", handlers.DeleteQuestion)

		protected.POST("/forms/:id/questions/:question_id/answers", handlers.CreateAnswer)
		protected.PUT("/forms/:id/questions/:question_id/answers/:answer_id", handlers.UpdateAnswer)
		protected.DELETE("/forms/:id/questions/:question_id/answers/:answer_id", handlers.DeleteAnswer)
	}
	if err := r.Run(forms_port); err != nil {
		panic(err)
	}
}
