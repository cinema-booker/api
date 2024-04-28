package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cinema-booker/api/api"
	"github.com/cinema-booker/api/config"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Error loading .env file: %v", err)
	}

	// connect to the database
	db, err := config.NewDatabase(config.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatalf("❌ Error connecting to the database: %v", err)
	}
	defer db.Close()

	// create & run an API server
	address := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	server := api.NewAPIServer(address, db)
	if err := server.Start(); err != nil {
		log.Fatalf("❌ Error starting the server: %v", err)
	}
}
