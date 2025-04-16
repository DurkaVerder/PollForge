package handlers

import (
	"auth/internal/models"
	"auth/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserLogging(c *gin.Context) {

	var request models.RegisterRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Printf("Ошибка ввода данных")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Инвалидный ввод"})
		return
	}

	err = service.CheckUserRequest(request)
	if err != nil {
		log.Printf("$1", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	token := service.GenerateJwt(request.id)
}
func UserRegistration(c *gin.Context) {

	var request models.RegisterRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Printf("invalid input")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	err = service.CheckUserRequest(request)
	token, err := service.RegistrUser(request)
	
}
