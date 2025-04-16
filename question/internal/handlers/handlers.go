package handlers

import (
	"question/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetAllQuestions() ([]models.Question, error)
	SubmitAnswer(answer models.Answer) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandlerAllQuestions(c *gin.Context) {

}

func (h *Handler) HandlerSubmitAnswer(c *gin.Context) {

}
