package main

import (
	"comments/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	router.SetUpRouter(r)
}
