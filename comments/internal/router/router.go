package router

import (
	"comments/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetUpRouter(r *gin.Engine) {

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/comments", handlers.GetComments)
	r.POST("/comments", handlers.CreateComment)
	r.PUT("/comments/:id", handlers.UpdateComment)
	r.DELETE("/comments/:id", handlers.DeleteComment)

}
