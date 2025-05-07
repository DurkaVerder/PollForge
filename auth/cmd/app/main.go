package main

import (
	"auth/internal/router"
	"auth/internal/storage"
	"auth/internal/kafka"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Print("Запуск микросервиса авторизации")
	err := storage.ConnectToDb()
	if err != nil {
		log.Fatal("Ошибка подключения к дб")
	}
	kafka.InitProducer()
	r := gin.Default()
	router.SetUpRouter(r)
}
