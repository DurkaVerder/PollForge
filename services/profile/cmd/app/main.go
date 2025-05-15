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
	os.MkdirAll("/uploads/avatars", 0755)
	r := gin.Default()
	r.Static("/avatars", "/uploads/avatars")
	r.MaxMultipartMemory = 8 << 20
	router.SetUpRouter(r)
}
