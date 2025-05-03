package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func extractFormID(c *gin.Context) (int, error) {
	formIdstr := c.Param("id")
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return 0, fmt.Errorf("неправильный тип id: %v", formIdstr)
	}
	return formId, nil
}

func GetComments(c *gin.Context) {
	
}

func CreateComment(c *gin.Context) {
	
}

func UpdateComment(c *gin.Context) {
	
}

func DeleteComment(c *gin.Context) {
	
}
