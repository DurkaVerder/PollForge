package websocket

import "github.com/gorilla/websocket"

type Service interface {
}

type WebSocket struct {
	upgrader *websocket.Upgrader

	clients map[string]*websocket.Conn

	service Service
}

func NewWebSocket(upgrader *websocket.Upgrader, service Service) *WebSocket {
	return &WebSocket{
		upgrader: upgrader,
		clients:  make(map[string]*websocket.Conn),
		service:  service,
	}
}
