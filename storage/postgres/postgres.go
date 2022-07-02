package postgres

import (
	"database/sql"
	_ "embed"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	db *sql.DB
}

//go:embed schema.sql
var queryCreate string

func NewStorage(databaseURI string) (*Storage, error) {
	db, err := sql.Open("pgx", databaseURI)
}
