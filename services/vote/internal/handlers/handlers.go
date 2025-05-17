package handlers

import (
	"log"
	"net/http"
	"question/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	AddVoteRequestToChannel(vote models.Vote)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandlerVote(c *gin.Context) {
	log.Println("HandlerVote: Получен запрос на голосование")
	var voteRequest models.VoteRequest
	if err := c.ShouldBindJSON(&voteRequest); err != nil {
		log.Println("HandlerVote: Ошибка при привязке JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userIDAny, exists := c.Get("userID")
	if !exists {
		log.Println("HandlerVote: Ошибка при получении userID из контекста")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	userID, ok := userIDAny.(int)
	if !ok {
		log.Println("HandlerVote: Ошибка при преобразовании userID в строку")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	vote := models.Vote{
		ID:       voteRequest.ID,
		IsUpVote: voteRequest.IsUpVote,
		UserID:   userID,
	}

	h.service.AddVoteRequestToChannel(vote)

	c.JSON(http.StatusOK, gin.H{"message": "Vote received"})
}
