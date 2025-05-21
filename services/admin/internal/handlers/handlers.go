package handlers

import (
	"admin/internal/models"
	"admin/internal/service"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось получить пользователей в сервисе админа"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func ToggleBanUser(c *gin.Context) {
	idstr := c.Param("id")

	if idstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не указан"})
		return
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id пользователя"})
		return
	}

	var ToogleBan models.ToogleBan
	if err := c.ShouldBindJSON(&ToogleBan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Ошибка при получении данных блокировки пользователя"})
		return
	}
	err = service.ToggleBanUser(id, ToogleBan.IsBanned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось изменить статус блокировки пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Сообщение": "Статус блокировки пользователя изменён"})
}

func DeleteUser(c *gin.Context) {
	idstr := c.Param("id")
	if idstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не указан"})
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id пользователя"})
		return
	}
	err = service.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Сообщение": "Пользователь успешно удалён"})
}


func DeleteForm(c *gin.Context) {
	formIdstr := c.Param("id")
	if formIdstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id формы не указан"})
		return
	}
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}
	err = service.FormDelete(formId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить форму"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно удалена"})
}


func DeleteComment(c *gin.Context) {
	commentIdstr := c.Param("id")
	if commentIdstr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id комментария не указан"})
		return
	}
	commentId, err := strconv.Atoi(commentIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id комментария"})
		return
	}
	err = service.DeleteComment(commentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить комментарий"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Сообщение": "Комментарий успешно удалён"})
}