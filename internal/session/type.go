package session

import (
	"time"

	"github.com/cinema-booker/internal/room"
)

type Session struct {
	Id        int        `json:"id" db:"id"`
	Price     int        `json:"price" db:"price"`
	StartsAt  time.Time  `json:"starts_at" db:"starts_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	Room      room.Room  `json:"room"`
}
