package kafka

import (
	"auth/internal/models"
	"context"
	"encoding/json"
	"log"
	"os"


	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer
var MsgChan chan models.MessageKafka

func InitProducer() {
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{os.Getenv("KAFKA_BROKER")},
		Topic:        `notify.topic`,
		Balancer:     &kafka.LeastBytes{},
	})

	ch := make(chan models.MessageKafka, 1000)
	MsgChan = ch
	go func() {
		for msg := range ch {
			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Ошибка маршалинга json: %v", err)
				continue
			}
			if err := writer.WriteMessages(context.Background(), kafka.Message{Value: data}); err != nil {
				log.Printf("Ошибка при отправке kafka-сообщения: %v", err)
			}
		}
	}()

}

func SendMessage(msg models.MessageKafka) error {
	select {
	case MsgChan <- msg:
	default:
		log.Println("MsgChan is full, dropping message")
	}
	return nil
}

