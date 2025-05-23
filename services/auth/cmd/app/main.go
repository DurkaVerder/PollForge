package main

import (
	"auth/internal/kafka"
	"auth/internal/router"
	"auth/internal/service"
	"auth/internal/storage"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
)

func main() {

	log.Print("Запуск микросервиса авторизации")
	err := storage.ConnectToDb()
	if err != nil {
		log.Fatal("Ошибка подключения к дб")
	}
	kafka.InitProducer()
	numWorkers := runtime.NumCPU()

	service.StartWorkerPool(numWorkers * 10)

	r := gin.Default()
	router.SetUpRouter(r)
}
