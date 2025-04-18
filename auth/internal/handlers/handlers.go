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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный ввод"})
		return
	}

	token, err := service.LoggingUser(request)
	if err != nil {
		log.Printf("Ошибка при входе пользователя в аккаунт")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ввод"})
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
func GetProfile(c *gin.Context) {
	id, exist := c.Get("id")
	if !exist {
		c.JSON(http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}
	profile, err := service.GetUserProfile(id.(int))

	if err != nil{
		c.JSON(http.StatusNotFound, "Ошибка при получении профиля")
		return
	}
	c.JSON(http.StatusOK,profile)

}
func GetForms(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	forms, err := service.GetUserForms(id.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении форм"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"forms": forms})
}
