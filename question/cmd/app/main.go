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

	service := service.NewService(postgres, answersChannel)

	wg := &sync.WaitGroup{}
	service.StartWorker(countWorkers, wg)

	handlers := handlers.NewHandler(service)

	engine := gin.Default()
	engine.Use(gin.Logger())

	ctx, cancel := context.WithCancel(context.Background())

	server := server.NewServer(handlers, engine)
	go func() {
		server.Start(os.Getenv("PORT"))
		cancel()
	}()

	<-ctx.Done()

	service.Close()
	wg.Wait()
}
