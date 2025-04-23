package kafka

import (
	"email/internal/models"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type KafkaEmailHandler struct {
	logger *log.Logger

	msgChan chan<- models.MessageKafka
}

func NewKafkaEmailHandler(logger *log.Logger, msgChan chan<- models.MessageKafka) *KafkaEmailHandler {
	return &KafkaEmailHandler{
		logger:  logger,
		msgChan: msgChan,
	}
}

func (h *KafkaEmailHandler) Setup(sess sarama.ConsumerGroupSession) error {
	h.logger.Println("Kafka consumer group setup")
	return nil
}
func (h *KafkaEmailHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	h.logger.Println("Kafka consumer group cleanup")
	return nil
}

func (h *KafkaEmailHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		var message models.MessageKafka
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			h.logger.Printf("Error unmarshaling message: %v", err)
			continue
		}

		h.msgChan <- message
		h.logger.Printf("Message sent to channel: %s", message)

		sess.MarkMessage(msg, "")
	}
	return nil
}
