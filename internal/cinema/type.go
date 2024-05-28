package cinema

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	RoomTypeSmall  = "SMALL"
	RoomTypeMedium = "MEDIUM"
	RoomTypeLarge  = "LARGE"
)

type Address struct {
	Id        int     `json:"id" db:"id"`
	Country   string  `json:"country" db:"country"`
	City      string  `json:"city" db:"city"`
	ZipCode   string  `json:"zip_code" db:"zip_code"`
	Street    string  `json:"street" db:"street"`
	Longitude float64 `json:"longitude" db:"longitude"`
	Latitude  float64 `json:"latitude" db:"latitude"`
}

type Room struct {
	Id     int    `json:"id" db:"id"`
	Number string `json:"number" db:"number"`
	Type   string `json:"type" db:"type"`
}

type Cinema struct {
	Id          int          `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	DeletedAt   sql.NullTime `json:"deleted_at" db:"deleted_at"`
	Address     Address      `json:"address"`
}

type CinemaWithRooms struct {
	Id          int          `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	DeletedAt   sql.NullTime `json:"deleted_at" db:"deleted_at"`
	Address     Address      `json:"address"`
	Rooms       RoomArray    `json:"rooms"`
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
