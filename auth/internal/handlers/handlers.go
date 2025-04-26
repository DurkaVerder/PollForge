package handlers

import (
	"auth/internal/models"
	"auth/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserLogging(c *gin.Context) {

	var request models.UserRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Printf("Ошибка ввода данных")
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный ввод"})
		return
	}

	token, err := service.LoggingUser(request)
	if err != nil {
		log.Printf("Ошибка при входе пользователя в аккаунт")
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь вошёл",
		"token":   token,
	})
}
func UserRegistration(c *gin.Context) {

	var request models.UserRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Printf("invalid input")
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Некорректный ввод"})
		return
	}
	err = service.CheckUserRequest(request)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие похожего аккаунта")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	token, err := service.RegisterUser(request)
	if err != nil {
		log.Printf("Ошибка при регистрации аккаунта")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Пользователь создан",
		"token":   token,
	})
}
