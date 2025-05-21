package main

import (
	"admin/internal/router"
	"log"
	"admin/internal/storage"
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
