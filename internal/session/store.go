package session

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SessionStore interface {
	FindById(id int) (Session, error)
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

func (s *Store) FindById(id int) (Session, error) {
	session := Session{}
	query := `
    SELECT 
      s.id AS id,
      s.price AS price,
      s.starts_at AS starts_at,
			s.deleted_at AS deleted_at,
      r.id AS "room.id",
      r.number AS "room.number",
      r.type AS "room.type"
    FROM sessions s
    LEFT JOIN rooms r ON s.room_id = r.id
		WHERE s.id=$1
  `
	err := s.db.Get(&session, query, id)

	return session, err
}

func (s *Store) Create(input map[string]interface{}) error {
	query := "INSERT INTO sessions (event_id, room_id, price, starts_at) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, input["event_id"], input["room_id"], input["price"], input["starts_at"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	columns := make([]string, 0, len(input))
	values := make([]interface{}, 0, len(input))
	for col, val := range input {
		columns = append(columns, fmt.Sprintf("%s=$%d", sqlx.Rebind(sqlx.DOLLAR, col), len(columns)+1))
		values = append(values, val)
	}

	query := fmt.Sprintf("UPDATE sessions SET %s WHERE id=$%d", strings.Join(columns, ", "), len(columns)+1)
	_, err := s.db.Exec(query, append(values, id)...)

	return err
}
