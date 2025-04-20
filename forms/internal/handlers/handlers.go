package handlers

import (
	"database/sql"
	"forms/internal/models"
	"forms/internal/service"
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

	// Конвертируем в численный тип данных строку с id для проверки
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
	// Проверка на наличие возвращаемого значения, если форм нет - сработает условие цикла
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

	forms, err := service.FormsGet(creatorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось найти анкеты"})
		return
	}
	c.JSON(http.StatusOK, forms)
}

func UpdateForm(c *gin.Context) {
	formIdstr := c.Param("id")

	// Конвертируем в численный тип данных строку с id для проверки
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
	// Конвертируем в численный тип данных строку с id для проверки
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

	// Проверка на существование формы для удаления, нужен id пользователя и id формы
	err = service.FormChek(creatorId, formId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	// Удаление формы, для этого требуется id формы и пользователя из-за нужды нахождения формы для удаления с помощью sql-запроса
	_, err = service.FormDelete(formId, creatorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно удалена"})
}
