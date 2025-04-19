package handlers

import (
	"forms/internal/models"
	"forms/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateForm(c *gin.Context) {
	var form models.FormRequest

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат данных"})
		return
	}

	creatorId, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пользователя не найден"})
		return
	}

	link := uuid.New().String()
	query := `INSERT INTO forms (creator_id, title, description, link, private_key, expires_at) 
			  VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	var formId int
	err = storage.Db.QueryRow(query, creatorId, form.Title, form.Description, link, form.PrivateKey, form.ExpiresAt).Scan(&formId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"form_id": formId, "link": link})
}

func GetForm(c *gin.Context) {

	creatorId, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пользователя не найден"})
		return
	}
	query := `SELECT id, title, description, link, private_key, expires_at FROM forms WHERE creator_id = $1`
	var form models.Form
	err := storage.Db.QueryRow(query, creatorId).Scan(&form.Id,
		&form.Title,
		&form.Description,
		&form.Link,
		&form.PrivateKey,
		&form.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось найти форму"})
		return
	}
}

func GetForms(c *gin.Context) {
	creatorId, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id пользователя не найден"})
		return
	}
	query := `SELECT id, title, description, link, private_key, expires_at FROM forms WHERE creator_id = $1`

	rows, err := storage.Db.Query(query, creatorId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось найти формы"})
		return
	}
	var forms []models.Form
	for rows.Next() {
		var form models.Form
		err := rows.Scan(&form.Id,
			&form.Title,
			&form.Description,
			&form.Link,
			&form.PrivateKey,
			&form.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось считать данные формы"})
			return
		}
		forms = append(forms, form)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, forms)
}

func UpdateForm(c *gin.Context) {

}

func DeleteForm(c *gin.Context) {

}
