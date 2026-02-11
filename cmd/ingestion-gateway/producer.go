package main

import (
	"context"
	"log"

	"github.com/JacksonDahl2/drone-platform/cmd/shared"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	producer *kafka.Writer
	topic    string
}

// function to create a new producer -> a constructor
// notice it's NOT part of the actual struct, unlike the following ones
func NewKafkaProducer(topic string) *KafkaProducer {
	cfg := shared.NewKafkaConfig()

	if topic == "" {
		topic = cfg.Topic
	}

	p := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{cfg.Host},
		Topic:   topic,
	})
	p.Completion = func(messages []kafka.Message, err error) {
		if err != nil {
			log.Printf("producer error: %v", err)
			return
		}
		for _, msg := range messages {
			log.Printf("delivered: topic=%s partition=%d offset=%d", msg.Topic, msg.Partition, msg.Offset)
		}
	}

	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}

// this will be the function call that actually produces the message
func (p *KafkaProducer) Produce(msg string) error {
	ctx := context.Background()

	return p.producer.WriteMessages(ctx, kafka.Message{Value: []byte(msg)})
}
