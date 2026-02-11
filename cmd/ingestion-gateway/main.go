package main

import (
	"io"
	"log"
	"net/http"
)

// logs all requests
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type Server struct {
	producer *KafkaProducer
	mux      *http.ServeMux
}

func NewServer(p *KafkaProducer) *Server {
	s := &Server{
		producer: p,
		mux:      http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("POST /ingest", s.handleIngest) // will be one for now, but expand to different types eventually
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("healthy")

	w.Write([]byte("Server healthy"))
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleIngest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer r.Body.Close()
	if err := s.producer.Produce(string(body)); err != nil {
		http.Error(w, "produce failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func main() {
	producer := NewKafkaProducer("")
	defer producer.Close()

	server := NewServer(producer)

	log.Print("Server starting on port 3000....")
	if err := http.ListenAndServe(":3000", logging(server)); err != nil {
		log.Fatal("Server crashed, ", err)

		return
	}
}
