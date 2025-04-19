package main

import (
	"net/http"
	"os"
	"stats/internal/server"
	"stats/internal/service"
	"stats/internal/storage"
	wb "stats/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	driverName = "postgres"
)

func main() {
	db, err := storage.InitConnection(driverName, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	postgres := storage.NewPostgres(db)

	svc := service.NewService(postgres)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	websocket := wb.NewWebSocket(&upgrader, svc)

	engine := gin.Default()

	srv := server.NewServer(websocket, engine)

	srv.Start(os.Getenv("PORT"))

}
