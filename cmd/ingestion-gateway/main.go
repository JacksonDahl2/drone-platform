package main

import (
	"net/http"

	"github.com/JacksonDahl2/drone-platform/cmd/shared"
)

type Server struct {
	gpsProducer    *KafkaProducer
	stateProducer  *KafkaProducer
	eventsProducer *KafkaProducer
	mux            *http.ServeMux
}

const (
	TopicGPS    = "v1_gps"
	TopicState  = "v1_state"
	TopicEvents = "v1_events"
)

func NewServer() *Server {
	s := &Server{
		gpsProducer:    NewKafkaProducer(TopicGPS),
		stateProducer:  NewKafkaProducer(TopicState),
		eventsProducer: NewKafkaProducer(TopicEvents),
		mux:            http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) Close() {
	_ = s.gpsProducer.Close()
	_ = s.stateProducer.Close()
	_ = s.eventsProducer.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func main() {
	server := NewServer()
	defer server.Close()
	shared.RunServer(":3000", "ingestion gateway", server)
}
