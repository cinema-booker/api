package event

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EventStore interface {
	FindAll() ([]Event, error)
	FindById(id int) (Event, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) FindAll() ([]Event, error) {
	events := []Event{}
	query := "SELECT * FROM events"
	err := s.db.Select(&events, query)

	return events, err
}

func (s *Store) FindById(id int) (Event, error) {
	event := Event{}
	query := "SELECT * FROM events WHERE id=$1"
	err := s.db.Get(&event, query, id)

	return event, err
}

func (s *Store) Create(input map[string]interface{}) error {
	query := "INSERT INTO events (room_id, price, starts_at, ends_at) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, input["room_id"], input["price"], input["starts_at"], input["ends_at"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	columns := make([]string, 0, len(input))
	values := make([]interface{}, 0, len(input))
	for col, val := range input {
		columns = append(columns, fmt.Sprintf("%s=$%d", sqlx.Rebind(sqlx.DOLLAR, col), len(columns)+1))
		values = append(values, val)
	}

	query := fmt.Sprintf("UPDATE events SET %s WHERE id=$%d", strings.Join(columns, ", "), len(columns)+1)
	_, err := s.db.Exec(query, append(values, id)...)

	return err
}
