package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func NewDatabase(config DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		config.Host, config.Port, config.User, config.Password, config.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	log.Println("🚀 Successfully connected to the database.")
	return db, nil
}
