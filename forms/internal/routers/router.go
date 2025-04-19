package router

import (
	"forms/internal/handlers"
	"forms/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {

	forms_port := os.Getenv("PORT")
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/forms",     handlers.CreateForm)
		protected.GET("/forms/:id",  handlers.GetForm)
		protected.PUT("/forms/:id",  handlers.UpdateForm)
		protected.DELETE("/forms/:id", handlers.DeleteForm)
	}
	if err := r.Run(forms_port); err != nil {
		panic(err)
	}
}


