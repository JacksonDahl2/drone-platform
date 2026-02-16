package main

import (
	"context"
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

func (*Processor) ProcessGps(ctx context.Context, msg []byte) error {
	return nil
}

func (*Processor) ProcessState(ctx context.Context, msg []byte) error {
	return nil
}

func (*Processor) ProcessEvent(ctx context.Context, msg []byte) error{
	return nil
}