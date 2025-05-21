package server

import (
	"stream_line/internal/handlers"
	"stream_line/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
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

	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://pollforge.ru"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	streamLine := s.engine.Group("/api/streamline")

	streamLine.Use(middleware.JWTAuth())
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
