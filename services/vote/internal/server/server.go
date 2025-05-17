package server

import (
	"log"
	"question/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	userIDKey = "userID"
)

type Handlers interface {
	HandlerVote(c *gin.Context)
}

type Server struct {
	engine   *gin.Engine
	handlers Handlers
}

func NewServer(handlers Handlers, engine *gin.Engine) *Server {
	return &Server{
		engine:   engine,
		handlers: handlers,
	}
}

func (s *Server) initRoutes() {
	s.engine.Use(Logger())

	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	vote := s.engine.Group("/vote")
	vote.Use(Authorization())
	vote.POST("/input", s.handlers.HandlerVote)
}

func (s *Server) Start(port string) {
	s.initRoutes()

	if err := s.engine.Run(port); err != nil {
		panic(err)
	}
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		log.Printf("Request: %s %s", ctx.Request.Method, ctx.Request.URL.Path)

		ctx.Next()

		latency := time.Since(t)
		log.Printf("Response: %d %s in %v", ctx.Writer.Status(), ctx.Request.URL.Path, latency)
	}
}

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		userID, err := service.GetParamFromJWT(authHeader, userIDKey)
		if err != nil {
			ctx.JSON(401, gin.H{"error": "Authorization header is required", "message": err.Error()})
			ctx.Abort()
			return
		}

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			ctx.JSON(401, gin.H{"error": "Invalid user ID", "message": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userIDInt)

		ctx.Next()
	}
}
