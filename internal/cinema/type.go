package cinema

import (
	"time"

	"github.com/cinema-booker/internal/room"
)

type Address struct {
	Id        int     `json:"id" db:"id"`
	Address   string  `json:"address" db:"address"`
	Longitude float64 `json:"longitude" db:"longitude"`
	Latitude  float64 `json:"latitude" db:"latitude"`
}

type Cinema struct {
	Id          int        `json:"id" db:"id"`
	UserId      int        `json:"user_id" db:"user_id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
	Address     Address    `json:"address"`
}

type CinemaWithRooms struct {
	Id          int            `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description string         `json:"description" db:"description"`
	DeletedAt   *time.Time     `json:"deleted_at" db:"deleted_at"`
	Address     Address        `json:"address"`
	Rooms       room.RoomArray `json:"rooms"`
}
