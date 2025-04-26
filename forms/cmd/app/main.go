package main

import (
	"forms/internal/router"
	"forms/internal/storage"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки из env. файла")
	}

}

func main() {

	err := storage.ConnectToDb()
	if err != nil {
		log.Fatal("Ошибка подключения к дб")
	}
	r := gin.Default()
	router.SetUpRouter(r)
}
