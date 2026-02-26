package main

import (
	"database/sql"

	sqlc "github.com/JacksonDahl2/drone-platform/internal/platform/db/sqlc"
	

)

type Accessor struct {
	queries *sqlc.Queries
}

func NewAccessor(db *sql.DB) *Accessor {
	queries := sqlc.New(db)
	return &Accessor{
		queries: queries,
	}
}