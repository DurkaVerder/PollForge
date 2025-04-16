package router

import (
	"auth/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine){
	r.POST("/register", handlers.UserRegistration)
	r.POST("/logging", handlers.UserLogging)
	r.Run(":8081")
}