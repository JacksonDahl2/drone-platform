package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JacksonDahl2/drone-platform/cmd/shared/models"
)

// POST /gps (coords, altitude, orientation, velocity, angular rate)
// POST /state (status, battery, connection)
// POST /events (mission started/ended, waypoint reached, alerts about battery, lost link, error)

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)

	s.mux.HandleFunc("POST /gps", s.handleGps)
	s.mux.HandleFunc("POST /state", s.handleState)
	s.mux.HandleFunc("POST /events", s.handleEvents)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("healthy")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server healthy"))
}

func (s *Server) handleGps(w http.ResponseWriter, r *http.Request) {
	log.Printf("ingestion-gateway.handleGps - hit")
	var payload models.GpsInput
	// decode
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("ingestion-gateway.handleGps - invalid input")
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// convert to bytes
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// send to producer
	if err := s.gpsProducer.Produce(string(data)); err != nil {
		log.Printf("produce failed: %v", err)
		http.Error(w, "produce failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Successfully submitted"))
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	log.Printf("ingestion-gateway.handleState - hit")
	var payload models.StateInput
	// decode
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// convert to bytes
	data, _ := json.Marshal(payload)

	// send to producer
	if err := s.stateProducer.Produce(string(data)); err != nil {
		http.Error(w, "produce failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Successfully submitted"))
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	log.Printf("ingestion-gateway.handleEvents - hit")
	var payload models.EventInput
	// decode
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// convert to bytes
	data, _ := json.Marshal(payload)

	// send to producer
	if err := s.eventsProducer.Produce(string(data)); err != nil {
		http.Error(w, "produce failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Successfully submitted"))
}

// func (s *Server) handleIngest(w http.ResponseWriter, r *http.Request) 	{
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 	}
// 	defer r.Body.Close()
// 	if err := s.producer.Produce(string(body)); err != nil {
// 		http.Error(w, "produce failed", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusAccepted)
// }

// input routes that I want

