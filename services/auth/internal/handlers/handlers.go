package handlers

import (
	"auth/internal/models"
	"auth/internal/service"
	"auth/internal/storage"
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
	
	token, err := service.AsyncLoginUser(request)
	if err != nil {
		log.Printf("Ошибка при входе пользователя в аккаунт")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := storage.CheckingLoggingData(request)
	if err != nil {
		log.Printf("Ошибка при сопоставлении пароля и почты - loginUserInternal")
		log.Printf("%s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      userId,
		"message": "Пользователь вошёл",
		"token":   token,
	})
}
func UserRegistration(c *gin.Context) {

	var request models.UserRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Printf("Ошибка ввода данных")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ввод"})
		return
	}
	_, err = service.AsyncRegisterUser(request)
	if err != nil {
		log.Printf("Ошибка при регистрации аккаунта")
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь создан",
	})
}

func PasswordResetRequest(c *gin.Context) {
    var req models.PasswordResetRequest 
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ввод"})
        return
    }
    _, err := service.AsyncResetPassword(models.UserRequest{Email: req.Email})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    // возвращаем сообщение о отправленной ссылке
    c.JSON(http.StatusOK, gin.H{"message": "ссылка для сброса отправлена"})
}

func PasswordResetConfirmHandler(c *gin.Context) {
    var req models.PasswordResetConfirm  
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "некорректные данные"})
        return
    }
    if err := service.AsyncConfirmReset(req.Token, req.NewPassword); err != nil {
        log.Printf("Confirm reset error: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "пароль успешно изменён"})
}

