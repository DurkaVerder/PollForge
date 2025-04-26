package server

import (
	"stats/internal/websocket"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ws *websocket.WebSocket

	engine *gin.Engine
}

func NewServer(ws *websocket.WebSocket, engine *gin.Engine) *Server {
	return &Server{
		ws:     ws,
		engine: engine,
	}
}

func (s *Server) initRoutes() {
	s.engine.GET("/ws", s.ws.HandleConnection)
}

func (s *Server) Start(port string) {
	if port == "" {
		port = ":8080"
	}

	s.initRoutes()

	if err := s.engine.Run(port); err != nil {
		panic(err)
	}
}
