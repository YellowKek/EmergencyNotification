package internal

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

func NewDB() (*sql.DB, error) {
	connCfg, err := pgx.ParseURI("postgres://postgres:postgres@postgres:5432/EmergencyNotification")
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres connection string: %w", err)
	}
	db := stdlib.OpenDB(connCfg)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	return db, nil
}
