package router

import (
	"github.com/gin-gonic/gin"
	"comments/internal/handlers")

func SetUpRouter(r *gin.Engine) {

	r.GET("/comments", handlers.GetComments)
	r.POST("/comments", handlers.CreateComment)
	r.PUT("/comments/:id", handlers.UpdateComment)
	r.DELETE("/comments/:id", handlers.DeleteComment)

}