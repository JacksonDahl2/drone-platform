package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/JacksonDahl2/drone-platform/cmd/shared"
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

type MessageHandler func(ctx context.Context, msg []byte) error

func (r *KafkaConsumer) Consume(ctx context.Context, handler MessageHandler, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		m, err := r.consumer.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Printf("consumer %s: shutting down", r.topic)
			} else {
				log.Printf("consumer %s read error: %v", r.topic, err)
			}
			return
		}
		log.Printf("topic %s | message at offset %d: %s = %s\n", r.topic, m.Offset, string(m.Key), string(m.Value))
		if err := handler(ctx, m.Value); err != nil {
			log.Printf("consumer %s handle message: %v", r.topic, err)
		}
	}
}

func main() {
	log.Printf("starting consumer worker...")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://drone:drone@localhost:5432/drone_platform?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	processor := NewProcessor(db)

	cGps := NewKafkaConsumer("v1_gps", "consume-group")
	cState := NewKafkaConsumer("v1_state", "consume-group")
	cEvents := NewKafkaConsumer("v1_events", "consume-group")
	defer cGps.Close()
	defer cState.Close()
	defer cEvents.Close()

	var wg sync.WaitGroup
	wg.Add(3)
	go cGps.Consume(ctx, func(ctx context.Context, msg []byte) error {
		return processor.ProcessGps(ctx, msg)
	}, &wg)
	go cState.Consume(ctx, func(ctx context.Context, msg []byte) error {
		return processor.ProcessState(ctx, msg)
	}, &wg)
	go cEvents.Consume(ctx, func(ctx context.Context, msg []byte) error {
		return processor.ProcessEvent(ctx, msg)
	}, &wg)

	// Cleanup steps
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Printf("received shutdown signal, draining consumers...")
	cancel()
	wg.Wait()
	log.Printf("consumer worker stopped")
}
