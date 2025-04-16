package main

import (
	"context"
	"os"
	"question/internal/handlers"
	"question/internal/server"
	"question/internal/service"
	"question/internal/storage"
	"question/models"
	"sync"

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

	answersChannel := make(chan []models.SubmitAnswer, 100)

	questionService := service.NewService(postgres, answersChannel)

	wg := &sync.WaitGroup{}
	questionService.StartWorker(countWorkers, wg)

	handlers := handlers.NewHandler(questionService)

	engine := gin.Default()
	engine.Use(gin.Logger())

	ctx, cancel := context.WithCancel(context.Background())

	srv := server.NewServer(handlers, engine)
	go func() {
		srv.Start(os.Getenv("PORT"))
		cancel()
	}()

	<-ctx.Done()

	questionService.Close()
	wg.Wait()
}
