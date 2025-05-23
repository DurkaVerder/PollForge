package handler

import (
	"log"
	"net/http"
	"stats/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetProfileStats(userID string) (models.ProfileStatsRequest, error)
}

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) HandlerProfileStats(c *gin.Context) {
	userID := c.Param("user_id")

	request, err := h.s.GetProfileStats(userID)
	if err != nil {
		log.Printf("error getting profile stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, request)

}
