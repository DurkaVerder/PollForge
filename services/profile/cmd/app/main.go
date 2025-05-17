package main

import (
	"os"
	"profile/internal/router"
	"profile/internal/storage"

	"log"

	"github.com/gin-gonic/gin"
)


func main() {

	err := storage.ConnectToDb()
	if err != nil {
		log.Fatal("Ошибка подключения к дб")
	}
	err = os.MkdirAll("/uploads/avatars", 0755)
	if err != nil {
		log.Fatal("Ошибка создания директории для аватаров")
	}
	r := gin.Default()
	r.Static("/avatars", "/uploads/avatars")
	r.MaxMultipartMemory = 8 << 20
	router.SetUpRouter(r)
}
