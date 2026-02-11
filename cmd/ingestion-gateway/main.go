package main

import (
	"fmt"
	"log"
	"time"
)


type Server struct {
	producer *KafkaProducer
}

func NewServer() *Server {
	return &Server {
		producer: NewKafkaProducer(""),
	}
}

func (s *Server) produceMessage() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	id := 0
	for t := range ticker.C {
		msg := fmt.Sprintf("hello world, msgId =%d, ts =%s", id, t.Format("15:20:20"))
		s.producer.Produce(msg)
		id++
	}

}

func main() {
	log.Printf("Starting server ... ")
	s := NewServer()
	s.produceMessage()

}
