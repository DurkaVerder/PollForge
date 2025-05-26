package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"profile/internal/models"
	"profile/internal/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	NginxURL = "http://localhost:80"
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

	if profile.AvatarURL != "" {
		profile.AvatarURL = fmt.Sprintf("%s%s", NginxURL, profile.AvatarURL)
	}

	c.JSON(http.StatusOK, profile)

}
func GetForms(c *gin.Context) {
	userId, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, "id пользователя не найден")
		return
	}

	forms, err := service.GetUserForms(userId)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id формы"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пользователя не найден"})
		return
	}

	// Проверка на существование формы для удаления, нужен id пользователя и id формы
	err = service.FormCheck(creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на существование формы: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Форма не найдена"})
		return
	}

	// Удаление формы, для этого требуется id формы и пользователя из-за нужды нахождения формы для удаления с помощью sql-запроса
	_, err = service.FormDelete(formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при удалении формы: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Форма успешно удалена"})
}

func UpdateProfileName(c *gin.Context) {
	id, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id пользователя не найден"})
		return
	}

	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		log.Printf("Ошибка при получении данных профиля: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении данных профиля"})
		return
	}

	err = service.UpdateProfileName(id, profile)
	if err != nil {
		log.Printf("Ошибка при обновлении профиля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении профиля"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Профиль успешно обновлён"})
}

func DeleteProfile(c *gin.Context) {
	id, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id пользователя не найден"})
		return
	}

	err = service.DeleteProfile(id)
	if err != nil {
		log.Printf("Ошибка при удалении профиля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении профиля"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Профиль успешно удалён"})
}

func UploadAvatar(c *gin.Context) {
	id, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id пользователя не найден"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		log.Printf("Ошибка при получении файла: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении файла"})
		return
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".jpg") && !strings.HasSuffix(strings.ToLower(file.Filename), ".png") && !strings.HasSuffix(strings.ToLower(file.Filename), ".jpeg") {
		log.Printf("Ошибка: неподдерживаемый формат файла: %s", file.Filename)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поддерживаются только JPG и PNG"})
		return
	}
	openedFile, err := file.Open()
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке файла"})
		return
	}
	defer openedFile.Close()

	buffer := make([]byte, 512)
	if _, err := openedFile.Read(buffer); err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке файла"})
		return
	}

	mimeType := http.DetectContentType(buffer)
	if mimeType != "image/jpg" && mimeType != "image/png" && mimeType != "image/jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Поддерживаются только JPG и PNG"})
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", id, time.Now().Unix(), ext)
	filePath := fmt.Sprintf("/uploads/avatars/%s", filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		log.Printf("Ошибка при сохранении файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
		return
	}
	avatarURL := fmt.Sprintf("/avatars/%s", filename)

	err = service.UploadAvatar(id, avatarURL)
	if err != nil {
		log.Printf("Ошибка при загрузке аватара: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при загрузке аватара"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Аватар успешно загружен",
		"avatar_url": avatarURL,
	})
}

func UpdateProfileBio(c *gin.Context) {
	id, err := extractUserID(c)
	if err != nil {
		log.Printf("Ошибка при получении id пользователя: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id пользователя не найден"})
		return
	}

	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		log.Printf("Ошибка при получении данных профиля: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении данных профиля"})
		return
	}

	err = service.UpdateProfileBio(id, profile.Bio)
	if err != nil {
		log.Printf("Ошибка при обновлении профиля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении профиля"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Описание профиля успешно обновлено"})
}

func GetDifUserProfile(c *gin.Context) {

	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Printf("Ошибка при преобразовании id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный id пользователя"})
		return
	}

	profile, err := service.GetDifUserProfile(userId)
	if err != nil {
		log.Printf("Ошибка при получении профиля GetDifUserProfile: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Ошибка при получении профиля"})
		return
	}

	if profile.AvatarURL != "" {
		profile.AvatarURL = fmt.Sprintf("%s%s", NginxURL, profile.AvatarURL)
	}

	c.JSON(http.StatusOK, profile)
}
