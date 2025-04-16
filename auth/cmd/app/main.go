package main

import (
	"auth/internal/handlers"
	"auth/internal/storage"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){

	err := godotenv.Load()
	if err != nil{
		log.Fatal("Ошибка загрузки из env. файла")
	}
	err = storage.ConnectToDb()
	if err != nil{
		log.Fatal("Ошибка подключения к дб")
	}
	r := gin.Default()	
	
	r.POST("/register", handlers.UserRegistration)
	r.POST("/logging", handlers.UserLogging)
	r.Run(":8081")
}
