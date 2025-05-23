package handler

import (
	"stats/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) HandlerProfileStats(c *gin.Context) {

}
