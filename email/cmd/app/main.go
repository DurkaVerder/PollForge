package main

import (
	"context"
	en "email/internal/email_notifier"
	"email/internal/kafka"
	"email/internal/models"
	"email/internal/service"
	"email/internal/storage"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

const (
	driver   = "postgres"
	groupID  = "notify_group"
	topics   = "notify.topic"
	sizeChan = 1000
)

func main() {
	logger := log.New(os.Stdout, "kafka-email-notifier: ", log.LstdFlags|log.Lshortfile)

	db, err := storage.ConnectDB(driver, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	storage := storage.NewPostgres(db)

	emailNotifier := en.NewEmailNotifier(os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"))

	srv := service.NewService(storage, emailNotifier, logger)

	msgChan := make(chan models.MessageKafka, sizeChan)

	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	brokers := []string{os.Getenv("KAFKA_BROKER")}
	topics := []string{topics}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		panic(err)
	}
	defer consumerGroup.Close()

	handler := kafka.NewKafkaEmailHandler(logger, msgChan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		logger.Println("Received interrupt signal, shutting down...")
		cancel()
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			logger.Println("Error from consumer group:", err)
		}
	}()

	srv.StartWorker(msgChan)

	for {
		if err := consumerGroup.Consume(ctx, topics, handler); err != nil {
			logger.Println("Error consuming messages:", err)
			break
		}

		if ctx.Err() != nil {
			logger.Println("Context error:", ctx.Err())
			break
		}
	}

	logger.Println("Shutting down consumer group...")
	srv.StopWorker()
}
