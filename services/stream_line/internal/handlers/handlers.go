package handlers

import (
	"net/http"
	"stream_line/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetStreamLines(userID string) (*models.StreamLineResponse, error)
}

type StreamLineHandler struct {
	s Service
}

func NewStreamLineHandler(s Service) *StreamLineHandler {
	return &StreamLineHandler{
		s: s,
	}
}

func (h *StreamLineHandler) GetStreamLine(ctx *gin.Context) {
	id, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}

	userID, ok := id.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is not a string"})
		return
	}

	polls, err := h.s.GetStreamLines(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, polls)
}
