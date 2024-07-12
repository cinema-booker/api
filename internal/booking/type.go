package booking

import (
	"time"
)

type Room struct {
	Id     int    `json:"id" db:"id"`
	Number string `json:"number" db:"number"`
	Type   string `json:"type" db:"type"`
}

type Event struct {
	ID        int        `json:"id" db:"id"`
	Price     int        `json:"price" db:"price"`
	StartsAt  time.Time  `json:"starts_at" db:"starts_at"`
	EndsAt    time.Time  `json:"ends_at" db:"ends_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	Room      Room       `json:"room"`
}

type User struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Booking struct {
	Id         int        `json:"id" db:"id"`
	Place      string     `json:"place" db:"place"`
	CanceledAt *time.Time `json:"canceled_at" db:"canceled_at"`
	Event      Event      `json:"event"`
	User       User       `json:"user"`
}
