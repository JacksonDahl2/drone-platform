package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"

	// "time"

	"github.com/JacksonDahl2/drone-platform/cmd/shared"
	sqlc "github.com/JacksonDahl2/drone-platform/internal/platform/db/sqlc"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	consumer      *kafka.Reader
	topic         string
	consumerGroup string
}

func NewKafkaConsumer(topic string, consumerGroup string) *KafkaConsumer {
	cfg := shared.NewKafkaConfig()
	if topic == "" {
		topic = cfg.Topic
	}
	if consumerGroup == "" {
		consumerGroup = cfg.ConsumerGroup
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Host},
		Topic:    topic,
		GroupID:  consumerGroup,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		consumer:      r,
		topic:         topic,
		consumerGroup: consumerGroup,
	}
}

func (r *KafkaConsumer) Close() error {
	return r.consumer.Close()
}

func (r *KafkaConsumer) Consume() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	processor := NewProcessor(db)

	for {
		m, err := r.consumer.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("shutting down consumer")
			} else {
				log.Printf("consumer read error: %v", err)
			}
			break
		}
		log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

}

func main() {

	log.Printf("starting consumer worker...")
	c := NewKafkaConsumer("", "")
	c.Consume()
	defer c.Close()

}
