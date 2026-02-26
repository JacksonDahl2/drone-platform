package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/JacksonDahl2/drone-platform/cmd/shared"
	_ "github.com/lib/pq"
)

// This is the server that will serve the web clients -> routes that will serve the
// database -> as well as a route to set up a websocket connection

type Server struct {
	mux *http.ServeMux
	accessor *Accessor
}

func NewServer(db *sql.DB) *Server {
	accessor := NewAccessor(db)	
	s := &Server{
		mux: http.NewServeMux(),
		accessor: accessor,
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

	server := NewServer(db)
	defer server.Close()

	shared.RunServer(":3010", "client api", server)
}
