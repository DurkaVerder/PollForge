package handlers

import (
	"log"
	"net/http"
	"question/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetQuestions(formID string) (models.QuestionResponse, error)
	AddAnswerRequestToChannel(answer models.SubmitAnswerRequest)
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
	log.Println("HandlerAllQuestions: Получен запрос на получение всех вопросов")
	formID := c.Query("form_id")
	if formID == "" {
		log.Println("HandlerAllQuestions: form_id не указан")
		c.JSON(http.StatusBadRequest, gin.H{"error": "form_id is required"})
		return
	}

	questionsResponse, err := h.service.GetQuestions(formID)
	if err != nil {
		log.Printf("HandlerAllQuestions: Ошибка при получении вопросов: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	log.Printf("HandlerAllQuestions: Успешно получено %d вопросов\n", len(questionsResponse.Question))
	c.JSON(http.StatusOK, questionsResponse)
}

func (h *Handler) HandlerSubmitAnswer(c *gin.Context) {
	log.Println("HandlerSubmitAnswer: Получен запрос на отправку ответа")

	var request models.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("HandlerSubmitAnswer: Ошибка при разборе тела запроса: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("HandlerSubmitAnswer: Ответы для отправки: %+v\n", request.Answers)
	h.service.AddAnswerRequestToChannel(request)

	log.Println("HandlerSubmitAnswer: Ответ успешно отправлен")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
