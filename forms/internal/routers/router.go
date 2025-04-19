package router

import (
	"forms/internal/handlers"
	"forms/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/forms",     handlers.CreateForm)
		protected.GET("/forms/:id",  handlers.GetForm)
		protected.PUT("/forms/:id",  handlers.UpdateForm)
		protected.DELETE("/forms/:id", handlers.DeleteForm)
	}
	if err := r.Run(":8083"); err != nil {
		panic(err)
	}
}
