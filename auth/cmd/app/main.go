package main

import (
	"auth/internal/router"
	"auth/internal/storage"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Print("Запуск микросервиса авторизации")
	err := storage.ConnectToDb()
	if err != nil {
		log.Fatal("Ошибка подключения к дб")
	}
	r := gin.Default()
	router.SetUpRouter(r)
}
