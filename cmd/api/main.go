package api

import (
	"log"
	"net/http"

	"github.com/JacksonDahl2/drone-platform/cmd/shared"
)

// This is the server that will serve the web clients -> routes that will serve the
// database -> as well as a route to set up a websocket connection

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) Close() {
	log.Printf("server closing")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func main() {
	server := NewServer()
	defer server.Close()
	shared.RunServer(":3010", "client api", server)
}
