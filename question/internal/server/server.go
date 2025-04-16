package server

import "github.com/gin-gonic/gin"

type Handlers interface {
}

type Server struct {
	engine *gin.Engine

	handlers Handlers
}
