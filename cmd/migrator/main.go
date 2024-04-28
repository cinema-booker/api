package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
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
		fmt.Println("Migrate create with name", args[1])
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
		fmt.Println("Migrate up with step", step)
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
		fmt.Println("Migrate down with step", step)
	default:
		log.Fatal("unknown command")
	}
}
