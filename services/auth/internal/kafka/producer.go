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

func InitProducer() {
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    `notify.topic`,
		Balancer: &kafka.LeastBytes{},
	})
}

func SendMessage(msg models.MessageKafka) error {
	ctx := context.Background()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
	if err != nil {
		log.Printf("Ошибка при отправке Kafka-сообщения: %v", err)
		return err
	}
	log.Printf("Kafka-сообщение отправлено: %s", string(data))
	return nil
}
