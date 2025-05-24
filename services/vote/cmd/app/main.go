package main

import (
	"context"
	"os"
	"question/internal/handlers"
	"question/internal/models"
	"question/internal/server"
	"question/internal/service"
	"question/internal/storage"

	"github.com/gin-gonic/gin"
)

const (
	driverName   = "postgres"
	countWorkers = 5
)

func main() {
	db, err := storage.InitConnection(driverName, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	postgres := storage.NewPostgres(db)
	defer postgres.Close()

	answersChannel := make(chan models.Vote, 100)
	likesChannel := make(chan models.Like, 100)

	voteService := service.NewService(postgres, answersChannel, likesChannel)

	voteService.Start(countWorkers)

	handlers := handlers.NewHandler(voteService)

	engine := gin.Default()

	ctx, cancel := context.WithCancel(context.Background())

	srv := server.NewServer(handlers, engine)
	go func() {
		srv.Start(os.Getenv("PORT"))
		cancel()
	}()

	<-ctx.Done()

	voteService.Close()

	if err := postgres.Close(); err != nil {
		panic(err)
	}
}
