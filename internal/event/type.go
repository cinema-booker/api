package event

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cinema-booker/internal/cinema"
	"github.com/cinema-booker/internal/room"
)

type Movie struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Language    string `json:"language" db:"language"`
	Poster      string `json:"poster" db:"poster"`
	Backdrop    string `json:"backdrop" db:"backdrop"`
}

type SessionBasic struct {
	Id       int       `json:"id" db:"id"`
	Price    int       `json:"price" db:"price"`
	StartsAt time.Time `json:"starts_at" db:"starts_at"`
	Room     room.Room `json:"room"`
}

func (s *SessionBasic) UnmarshalJSON(data []byte) error {
	type Alias SessionBasic
	aux := &struct {
		StartsAt string `json:"starts_at"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, aux.StartsAt)
	if err != nil {
		// Try to parse without timezone
		parsedTime, err = time.Parse("2006-01-02T15:04:05", aux.StartsAt)
		if err != nil {
			return err
		}
	}
	s.StartsAt = parsedTime

	return nil
}

type SessionArray []SessionBasic

func (r *SessionArray) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, r)
	}
	return fmt.Errorf("unsupported data type: %T", src)
}

func (r SessionArray) Value() (driver.Value, error) {
	return json.Marshal(r)
}

type Event struct {
	Id        int           `json:"id" db:"id"`
	DeletedAt *time.Time    `json:"deleted_at" db:"deleted_at"`
	Cinema    cinema.Cinema `json:"cinema"`
	Movie     Movie         `json:"movie"`
	Sessions  SessionArray  `json:"sessions"`
}
