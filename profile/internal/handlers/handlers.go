package handlers

import (
	"fmt"
	"log"
	"net/http"
	"profile/internal/service"

	"github.com/gin-gonic/gin"
)

func extractUserID(c *gin.Context) (int, error) {
	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return 0, fmt.Errorf("id пользователя не найден")
	}
	creatorId, ok := creatorIdfl.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный тип id"})
		return 0, fmt.Errorf("неправильный тип id: %v", creatorIdfl)
	}
	return creatorId, nil
}

func GetProfile(c *gin.Context) {
	id, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, "id пользователя не найден")
		return
	}
	profile, err := service.GetUserProfile(id)

	if err != nil {
		log.Printf("Ошибка при получении профиля: %v", err)
		c.JSON(http.StatusNotFound, "Ошибка при получении профиля")
		return
	}
	c.JSON(http.StatusOK, profile)

}
func GetForms(c *gin.Context) {
	id, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, "id пользователя не найден")
		return
	}

	forms, err := service.GetUserForms(id)
	if err != nil {
		log.Printf("Ошибка при получении форм: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении форм"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"forms": forms})
}
