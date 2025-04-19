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
		protected.POST("/polls",     handlers.CreatePoll)
		protected.GET("/polls/:id",  handlers.GetPoll)
		protected.PUT("/polls/:id",  handlers.UpdatePoll)
		protected.DELETE("/polls/:id", handlers.DeletePoll)
	}
	if err := r.Run(":8083"); err != nil {
		panic(err)
	}
}
