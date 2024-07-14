package cinema

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Address struct {
	Id        int     `json:"id" db:"id"`
	Address   string  `json:"address" db:"address"`
	Longitude float64 `json:"longitude" db:"longitude"`
	Latitude  float64 `json:"latitude" db:"latitude"`
}

type Room struct {
	Id     int    `json:"id" db:"id"`
	Number string `json:"number" db:"number"`
	Type   string `json:"type" db:"type"`
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
	Id          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
	Address     Address    `json:"address"`
	Rooms       RoomArray  `json:"rooms"`
}

type RoomArray []Room

func (r *RoomArray) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, r)
	}
	return fmt.Errorf("unsupported data type: %T", src)
}

func (r RoomArray) Value() (driver.Value, error) {
	return json.Marshal(r)
}
