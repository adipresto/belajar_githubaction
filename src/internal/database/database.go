package database

import (
	"fmt"
	"log"
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func InitDBConnection(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Database, cfg.DB.SSL)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	return db, nil
}

func CloseDBConnection(db *sqlx.DB) error {
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
