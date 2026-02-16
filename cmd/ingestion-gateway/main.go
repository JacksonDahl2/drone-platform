package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// logs all requests
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// wraps the server to help it recover if there are any panics, so server doesn't crash
func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

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

	srv := &http.Server{
		Addr:    ":3000",
		Handler: recovery(logging(server)),
	}
	go func() {
		log.Print("Server starting on port 3000....")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	log.Printf("received shutdown signal, draining...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
	log.Printf("ingestion gateway stopped")
}
