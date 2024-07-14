package room

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Room struct {
	Id     int    `json:"id" db:"id"`
	Number string `json:"number" db:"number"`
	Type   string `json:"type" db:"type"`
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
