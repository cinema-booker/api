package room

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type RoomStore interface {
	FindById(id int) (Room, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	Delete(id int) error
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) FindById(id int) (Room, error) {
	room := Room{}
	query := "SELECT * FROM rooms WHERE id=$1"
	err := s.db.Get(&room, query, id)

	return room, err
}

func (s *Store) Create(input map[string]interface{}) error {
	query := "INSERT INTO rooms (cinema_id, number, type) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(query, input["cinema_id"], input["number"], input["type"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	columns := make([]string, 0, len(input))
	values := make([]interface{}, 0, len(input))
	for col, val := range input {
		columns = append(columns, fmt.Sprintf("%s=$%d", sqlx.Rebind(sqlx.DOLLAR, col), len(columns)+1))
		values = append(values, val)
	}

	query := fmt.Sprintf("UPDATE rooms SET %s WHERE id=$%d", strings.Join(columns, ", "), len(columns)+1)
	_, err := s.db.Exec(query, append(values, id)...)

	return err
}

func (s *Store) Delete(id int) error {
	query := "DELETE FROM rooms WHERE id=$1"
	_, err := s.db.Exec(query, id)

	return err
}
