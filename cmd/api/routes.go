package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// DB-backed endpoints (wire s.queries and parse params in handlers):
//
//	GET /api/drones                    -> GetLatestGpsAndStatePerDrone (or GetLatestGpsAllDrones + GetLatestStateAllDrones)
//	GET /api/drones/{id}/gps            -> GetGpsByDroneTimeRange  ?from=,&to= (RFC3339)
//	GET /api/drones/{id}/state          -> GetStateByDroneTimeRange ?from=,&to=
//	GET /api/drones/{id}/events         -> GetEventsByDroneTimeRange ?from=,&to=
//	GET /api/events                    -> GetEventsByTimeRange      ?from=,&to=
//	GET /api/events/recent             -> GetRecentEvents           ?limit=
//	GET /api/metrics/drone-count       -> GetDroneCount
//	GET /api/metrics/by-status         -> GetDroneCountByStatus
//	GET /api/metrics/connected        -> GetConnectedDroneCount
//	GET /api/metrics/events-by-type   -> GetEventCountByType       ?from=,&to=
//	GET /api/metrics/battery           -> GetBatteryStatsByDrone    ?from=,&to=
//	GET /api/metrics/activity         -> GetActivityByTimeBucket   ?bucket=,&from=,&to=  (e.g. bucket=5m)
//	GET /api/metrics/active-drones     -> GetActiveDronesPerTimeBucket ?bucket=,&from=,&to=
func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)

	s.mux.HandleFunc("GET /api/drones", s.handleLatestDrones)
	s.mux.HandleFunc("GET /api/drones/{id}/gps", s.handleDroneGpsRange)
	s.mux.HandleFunc("GET /api/drones/{id}/state", s.handleDroneStateRange)
	s.mux.HandleFunc("GET /api/drones/{id}/events", s.handleDroneEventsRange)

	s.mux.HandleFunc("GET /api/events", s.handleEventsTimeRange)
	s.mux.HandleFunc("GET /api/events/recent", s.handleRecentEvents)

	s.mux.HandleFunc("GET /api/metrics/drone-count", s.handleDroneCount)
	s.mux.HandleFunc("GET /api/metrics/by-status", s.handleDroneCountByStatus)
	s.mux.HandleFunc("GET /api/metrics/connected", s.handleConnectedCount)
	s.mux.HandleFunc("GET /api/metrics/events-by-type", s.handleEventCountByType)
	s.mux.HandleFunc("GET /api/metrics/battery", s.handleBatteryStats)
	s.mux.HandleFunc("GET /api/metrics/activity", s.handleActivityByBucket)
	s.mux.HandleFunc("GET /api/metrics/active-drones", s.handleActiveDronesPerBucket)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("server health"))
}

func (s *Server) handleLatestDrones(w http.ResponseWriter, r *http.Request) {
	log.Printf("routes.handleLatestDrones - entry")
	res, err := s.accessor.queries.GetLatestGpsAndStatePerDrone(r.Context())
	if err != nil {
		log.Printf("routes.handleLatestDrones - Failed query: %v", err)
		s.writeJSON(w, 500, []struct{}{})	
		return
	}

	s.writeJSON(w, http.StatusOK, res)
	return
}

func (s *Server) handleDroneGpsRange(w http.ResponseWriter, r *http.Request) {
	_ = r.PathValue("id")
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleDroneStateRange(w http.ResponseWriter, r *http.Request) {
	_ = r.PathValue("id")
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleDroneEventsRange(w http.ResponseWriter, r *http.Request) {
	_ = r.PathValue("id")
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleEventsTimeRange(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleRecentEvents(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleDroneCount(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]int64{"count": 0})
}

func (s *Server) handleDroneCountByStatus(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleConnectedCount(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]int64{"count": 0})
}

func (s *Server) handleEventCountByType(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleBatteryStats(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleActivityByBucket(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) handleActiveDronesPerBucket(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, []struct{}{})
}

func (s *Server) writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: %v", err)
	}
}