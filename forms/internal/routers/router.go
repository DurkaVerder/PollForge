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
		protected.POST("/polls",     handlers.CreateForm)
		protected.GET("/polls/:id",  handlers.GetForm)
		protected.PUT("/polls/:id",  handlers.UpdateForm)
		protected.DELETE("/polls/:id", handlers.DeleteForm)
	}
	if err := r.Run(":8083"); err != nil {
		panic(err)
	}
}
