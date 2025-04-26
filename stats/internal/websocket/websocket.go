package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"stats/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Service interface {
	GetPollStats(userID, formID string) (models.PollStatsResponse, error)
}

type WebSocket struct {
	upgrader *websocket.Upgrader

	service Service
}

func NewWebSocket(upgrader *websocket.Upgrader, service Service) *WebSocket {
	return &WebSocket{
		upgrader: upgrader,
		service:  service,
	}
}

func (ws *WebSocket) HandleConnection(c *gin.Context) {
	userID := c.Query("user_id")
	formID := c.Query("form_id")

	if userID == "" || formID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "userID and formID are required",
			"user_id": userID,
			"form_id": formID,
		})
		return
	}

	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				break
			}
		}
	}()

	log.Printf("Client %s connected", userID)

	go ws.readMessage(userID, formID, conn, done)
}

func (ws *WebSocket) readMessage(userID, formID string, conn *websocket.Conn, done <-chan struct{}) {

	ticker := time.NewTicker(2 * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			stats, err := ws.service.GetPollStats(userID, formID)
			if err != nil {
				log.Printf("Error getting poll stats for user %s: %v", userID, err)
				continue
			}

			data, err := json.Marshal(stats)
			if err != nil {
				log.Printf("Error marshalling poll stats for user %s: %v", userID, err)
				continue
			}

			ws.sendData(userID, conn, data)

		case <-done:
			log.Printf("Client %s disconnected", userID)
			conn.Close()
			return
		}
	}
}

func (ws *WebSocket) sendData(userID string, conn *websocket.Conn, message []byte) {

	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error writing message to client %s: %v", userID, err)
	}
}
