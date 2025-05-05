package router

import (
	"comments/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRouter(r *gin.Engine) {

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/forms/:form_id/comments", handlers.GetComments)
	r.POST("/forms/:form_id/comments", handlers.CreateComment)
	r.PUT("/forms/:form_id/comments/:id", handlers.UpdateComment)
	r.DELETE("/forms/:form_id/comments/:id", handlers.DeleteComment)

}
