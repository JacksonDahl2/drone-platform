package api

import (
	"log"
	"net/http"
)

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("health")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("server health"))
}