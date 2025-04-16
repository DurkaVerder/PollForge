package server

import "github.com/gin-gonic/gin"

type Handlers interface {
	HandlerAllQuestions(c *gin.Context)
	HandlerSubmitAnswer(c *gin.Context)
}

type Server struct {
	engine *gin.Engine
	handlers Handlers
}

func NewServer(handlers Handlers, engine *gin.Engine) *Server {
	return &Server{
		engine:   engine,
		handlers: handlers,
	}
}

func (s *Server) initRoutes() {
	s.engine.GET("/questions", s.handlers.HandlerAllQuestions)
	s.engine.POST("/submit_answer", s.handlers.HandlerSubmitAnswer)
}

func (s *Server) Start(port string) {
	s.initRoutes()

	if err := s.engine.Run(port); err != nil {
		panic(err)
	}
}
