package handlers

import (
	"net/http"
	"strconv"
	"stream_line/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 10
)

type Service interface {
	GetStreamLines(userID string, cursor time.Time, limit int) (*models.StreamLineResponse, error)
	GetPollByLink(userID, pollLink string) (*models.StreamLineResponse, error)
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

	limitStr := ctx.DefaultQuery("limit", "10")
	cursorStr := ctx.DefaultQuery("cursor", "")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	var cursor time.Time
	if cursorStr != "" {
		cursor, err = time.Parse("2006-01-02 15:04:05", cursorStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid cursor format"})
			return
		}
	} else {
		cursor = time.Now()
	}

	polls, err := h.s.GetStreamLines(userID, cursor, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if polls == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"polls":       []models.Polls{},
			"next_cursor": "",
			"hasMore":     false,
		})
		return
	}

	var nextCursor string
	if len(polls.Polls) > 0 {
		nextCursor = polls.Polls[len(polls.Polls)-1].CreatedAt
	}

	ctx.JSON(http.StatusOK, gin.H{
		"polls":       polls,
		"next_cursor": nextCursor,
		"hasMore":     len(polls.Polls) == limit,
	})
}

func (h *StreamLineHandler) GetPollByLink(ctx *gin.Context) {
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

	pollLink := ctx.Param("link")
	poll, err := h.s.GetPollByLink(userID, pollLink)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, poll)
}
