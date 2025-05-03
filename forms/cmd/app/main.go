package main

import (
	"forms/internal/router"
	"forms/internal/storage"

	"log"

	"github.com/gin-gonic/gin"
)



func main() {

	err := storage.ConnectToDb()
	if err != nil {
		log.Fatal("Ошибка подключения к дб")
	}
	r := gin.Default()
	router.SetUpRouter(r)
}
