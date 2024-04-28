package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cinema-booker/api/config"
	"github.com/cinema-booker/api/pkg/migrator"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Error loading .env file : %v", err)
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
		log.Fatalf("❌ Error connecting to the database : %v", err)
	}
	defer db.Close()

	// create a migrator instance
	m, err := migrator.NewMigrator(db, "/migrations")
	if err != nil {
		log.Fatalf("❌ Error creating migrator : %v", err)
	}

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("no command provided")
	}

	switch args[0] {
	case "create":
		if len(args) > 2 {
			log.Fatal("too many arguments")
		}
		if len(args) < 2 {
			log.Fatal("no migration name provided")
		}
		if err := m.CreateMigration(args[1]); err != nil {
			log.Fatalf("error creating migration: %v", err)
		}
	case "up":
		if len(args) > 2 {
			log.Fatal("too many arguments")
		}
		step := -1
		if len(args) == 2 {
			var err error
			step, err = strconv.Atoi(args[1])
			if err != nil {
				log.Fatalf("invalid step value: %v", args[1])
			}
		}
		fmt.Println("migrate up with step", step)
	case "down":
		if len(args) > 2 {
			log.Fatal("too many arguments")
		}
		step := -1
		if len(args) == 2 {
			var err error
			step, err = strconv.Atoi(args[1])
			if err != nil {
				log.Fatalf("invalid step value: %v", args[1])
			}
		}
		fmt.Println("migrate down with step", step)
	default:
		log.Fatal("unknown command")
	}
}
