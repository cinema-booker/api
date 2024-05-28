package booking

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type BookingStore interface {
	FindAll() ([]Booking, error)
	FindById(id int) (Booking, error)
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

func (s *Store) FindAll() ([]Booking, error) {
	bookings := []Booking{}
	query := "SELECT * FROM bookings"
	err := s.db.Select(&bookings, query)

	return bookings, err
}

func (s *Store) FindById(id int) (Booking, error) {
	booking := Booking{}
	query := "SELECT * FROM bookings WHERE id=$1"
	err := s.db.Get(&booking, query, id)

	return booking, err
}

func (s *Store) Create(input map[string]interface{}) error {
	query := "INSERT INTO bookings (user_id, room_id, place) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(query, input["user_id"], input["room_id"], input["place"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	columns := make([]string, 0, len(input))
	values := make([]interface{}, 0, len(input))
	for col, val := range input {
		columns = append(columns, fmt.Sprintf("%s=$%d", sqlx.Rebind(sqlx.DOLLAR, col), len(columns)+1))
		values = append(values, val)
	}

	query := fmt.Sprintf("UPDATE bookings SET %s WHERE id=$%d", strings.Join(columns, ", "), len(columns)+1)
	_, err := s.db.Exec(query, append(values, id)...)

	return err
}
