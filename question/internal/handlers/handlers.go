package handlers

import (
	"log"
	"net/http"
	"question/internal/service"
	"question/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetQuestions(formID string, userID any) (models.QuestionResponse, error)
	AddAnswerRequestToChannel(answer models.SubmitAnswerRequest)
	CreateAnsweredPolls(formID string, userID any) error
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

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("HandlerAllQuestions: Ошибка при получении userID из контекста")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	questionsResponse, err := h.service.GetQuestions(formID, userID)
	if err != nil {
		if err.Error() == service.ErrUserAlreadyAnswered {
			log.Println("HandlerAllQuestions: Пользователь уже ответил на вопрос")
			c.JSON(http.StatusForbidden, gin.H{"error": "User already answered"})
			return
		}

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

	formID := c.Query("form_id")
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("HandlerSubmitAnswer: Ошибка при получении userID из контекста")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	if err := h.service.CreateAnsweredPolls(formID, userID); err != nil {
		log.Printf("HandlerSubmitAnswer: Ошибка при создании записей о ответах: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	log.Println("HandlerSubmitAnswer: Ответ успешно отправлен")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
