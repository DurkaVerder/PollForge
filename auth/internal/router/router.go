package router

import (
	"auth/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(r *gin.Engine) {
	r.POST("/register", handlers.UserRegistration)
	r.POST("/logging", handlers.UserLogging)
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
