package handlers

import (
	"fmt"
	"log"
	"net/http"
	"profile/internal/models"
	"profile/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func extractUserID(c *gin.Context) (int, error) {
	creatorIdfl, ok := c.Get("id")
	if !ok {

		return 0, fmt.Errorf("id пользователя не найден")
	}
	creatorId, ok := creatorIdfl.(string)
	if !ok {
		return 0, fmt.Errorf("неправильный тип id: %v", ok)
	}
	return strconv.Atoi(creatorId)
}

func extractFormID(c *gin.Context) (int, error) {
	formIdstr := c.Param("id")
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		return 0, fmt.Errorf("неправильный тип id: %v", err.Error())
	}
	return formId, nil
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

func DeleteForm(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	// Проверка на существование формы для удаления, нужен id пользователя и id формы
	err = service.FormChek(creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на существование формы: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	// Удаление формы, для этого требуется id формы и пользователя из-за нужды нахождения формы для удаления с помощью sql-запроса
	_, err = service.FormDelete(formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при удалении формы: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно удалена"})
}

func UpdateProfileName(c *gin.Context) {
	id, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, "id пользователя не найден")
		return
	}

	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		log.Printf("Ошибка при получении данных профиля: %v", err)
		c.JSON(http.StatusBadRequest, "Ошибка при получении данных профиля")
		return
	}

	err = service.UpdateProfile(id, profile)
	if err != nil {
		log.Printf("Ошибка при обновлении профиля: %v", err)
		c.JSON(http.StatusInternalServerError, "Ошибка при обновлении профиля")
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Профиль успешно обновлён"})
}

func DeleteProfile(c *gin.Context) {
	id, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, "id пользователя не найден")
		return
	}

	err = service.DeleteProfile(id)
	if err != nil {
		log.Printf("Ошибка при удалении профиля: %v", err)
		c.JSON(http.StatusInternalServerError, "Ошибка при удалении профиля")
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Профиль успешно удалён"})
}