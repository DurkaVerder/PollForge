package handlers

import (
	"database/sql"
	"forms/internal/models"
	"forms/internal/service"
	"forms/internal/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateForm(c *gin.Context) {
	var form models.FormRequest

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorID, ok := creatorIdfl.(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}
	creatorId := int(creatorID)

	var formId int
	formId, link, err := service.FormCreate(form, creatorId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось создать форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"form_id": formId, "link": link})
}

func GetForm(c *gin.Context) {

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	formIdstr := c.Param("id")

	formId, err := strconv.Atoi(formIdstr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}
	var form models.Form

	creatorID, ok := creatorIdfl.(float64)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный тип id пользователя"})
        return
    }

	creatorId := int(creatorID)
	form, err = service.FormGet(creatorId, formId)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка при получении формы"})
		return
	}

	c.JSON(http.StatusOK, form)
}

func GetForms(c *gin.Context) {
	creatorIdfl, ok := c.Get("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	creatorID, ok := creatorIdfl.(float64)

    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный тип id пользователя"})
        return
    }

	creatorId := int(creatorID)

	query := `SELECT id, title, description, link, private_key, expires_at FROM forms WHERE creator_id = $1`

	rows, err := storage.Db.Query(query, creatorId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Не удалось найти формы"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось считать данные формы"})
			return
		}
		forms = append(forms, form)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, forms)
}

func UpdateForm(c *gin.Context) {
	formIdstr := c.Param("id")

	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorID, ok := creatorIdfl.(float64)

	var updateForm models.FormRequest

	err = c.ShouldBindJSON(&updateForm)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	creatorId := int(creatorID)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	err = service.FormChek(creatorId, formId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	err = service.FormUpdate(updateForm, creatorId, formId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка обновления данных"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно обновлена"})
}

func DeleteForm(c *gin.Context) {
	formIdstr := c.Param("id")

	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorIdfl, ok := c.Get("id")

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	creatorID, ok := creatorIdfl.(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}
	creatorId := int(creatorID)

	err = service.FormChek(creatorId, formId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}


	_, err = service.FormDelete(formId, creatorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно удалена"})
}
