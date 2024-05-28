package event

import (
	"time"
)

type Room struct {
	Id     int    `json:"id" db:"id"`
	Number string `json:"number" db:"number"`
	Type   string `json:"type" db:"type"`
}

type Movie struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Language    string `json:"language" db:"language"`
	Poster      string `json:"poster" db:"poster"`
	Backdrop    string `json:"backdrop" db:"backdrop"`
}

type Event struct {
	ID        int        `json:"id" db:"id"`
	Price     int        `json:"price" db:"price"`
	StartsAt  time.Time  `json:"starts_at" db:"starts_at"`
	EndsAt    time.Time  `json:"ends_at" db:"ends_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	Room      Room       `json:"room"`
	Movie     Movie      `json:"movie"`
}
