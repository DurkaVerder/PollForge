package handlers

import (
	"net/http"
	"profile/internal/service"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	id, exist := c.Get("id")
	if !exist {
		c.JSON(http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}
	profile, err := service.GetUserProfile(id.(int))

	if err != nil {
		c.JSON(http.StatusNotFound, "Ошибка при получении профиля")
		return
	}
	c.JSON(http.StatusOK, profile)

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
