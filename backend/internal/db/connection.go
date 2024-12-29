package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	return sqlx.Connect("postgres", databaseURL)
}
