package server

import (
	"stream_line/internal/handlers"
	"stream_line/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	handlers *handlers.StreamLineHandler

	engine *gin.Engine
}

func NewServer(handlers *handlers.StreamLineHandler, engine *gin.Engine) *Server {
	return &Server{
		handlers: handlers,
		engine:   engine,
	}
}

func (s *Server) initRoutes() {

	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	streamLine := s.engine.Group("/streamline")
	streamLine.Use(middleware.AuthMiddleware())
	streamLine.Use(middleware.LoggerMiddleware())
	{
		streamLine.GET("/news", s.handlers.GetStreamLine)
	}

}

func (s *Server) Start(port string) {

	s.initRoutes()

	if err := s.engine.Run(port); err != nil {
		panic(err)
	}
}
