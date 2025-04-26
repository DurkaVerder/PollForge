package main

import (
	"log"
	"os"
	"stream_line/internal/handlers"
	"stream_line/internal/server"
	"stream_line/internal/service"
	"stream_line/internal/storage"

	"github.com/gin-gonic/gin"
)

const (
	driver = "postgres"
)

func main() {

	logger := log.New(os.Stdout, "[StreamLine]", log.LstdFlags|log.Lshortfile)

	db, err := storage.ConnectDB(driver, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	postgres := storage.NewPostgres(db)

	src := service.NewService(postgres, logger)

	hand := handlers.NewStreamLineHandler(src)

	engine := gin.Default()

	srv := server.NewServer(hand, engine)

	srv.Start(os.Getenv("PORT"))
}
