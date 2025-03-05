package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pedroxer/auth-service/internal/config"
)

func ConnectToPG(cfg *config.Postgres) (*sql.DB, error) {
	DSN := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=%s",
		cfg.Db,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Sslmode,
	)

	db, _ := sql.Open("postgres", DSN)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ConnectToPG: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxConns)

	return db, nil
}
