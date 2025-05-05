package handlers

import (
	"comments/internal/models"
	"comments/internal/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func extractFormID(c *gin.Context) (int, error) {
	formIdstr := c.Param("form_id")
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return 0, fmt.Errorf("неправильный тип id: %v", formIdstr)
	}
	return formId, nil
}

func extractCommentID(c *gin.Context) (int, error) {
	Commentstr := c.Param("id")
	formId, err := strconv.Atoi(Commentstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id комментария"})
		return 0, fmt.Errorf("неправильный тип id: %v", Commentstr)
	}
	return formId, nil
}

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

func GetComments(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	comments, err := service.GetAllComments(formId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка получения комментариев"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func CreateComment(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный формат данных"})
		return
	}
	userId, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id пользователя"})
		return
	}
	err = service.CreateComment(comment, formId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка создания комментария"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Сообщение": "Комментарий успешно создан"})
}

func UpdateComment(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный формат данных"})
		return
	}
	userId, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id пользователя"})
		return
	}
	commentId, err := extractCommentID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id комментария"})
		return
	}

	err = service.UpdateUserComment(comment, commentId, formId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка обновления комментария"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Комментарий успешно обновлён"})
}

func DeleteComment(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}
	userId, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id пользователя"})
		return
	}
	commentId, err := extractCommentID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id комментария"})
		return
	}
	
	err = service.DeleteComment(commentId, formId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка удаления комментария"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Комментарий успешно удалён"})
}
