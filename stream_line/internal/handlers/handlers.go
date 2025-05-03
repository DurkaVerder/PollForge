package handlers

import (
	"net/http"
	"stream_line/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetStreamLine(userID string) (models.FormResponse, error)
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

	forms, err := h.s.GetStreamLine(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, forms)
}
