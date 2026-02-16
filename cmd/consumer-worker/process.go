package main

import (
	"database/sql"

	sqlc "github.com/JacksonDahl2/drone-platform/internal/platform/db/sqlc"
)

type Processor struct {
	queries  *sqlc.Queries
}

func NewProcessor(db *sql.DB) *Processor {
	queries := sqlc.New(db)
	return &Processor{
		queries: queries,
	}
}

func (*Processor) ProcessGps() {

}

func (*Processor) ProcessState() {

}

func (*Processor) ProcessEvent() {
	return
}