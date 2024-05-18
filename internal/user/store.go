package user

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserStore interface {
	FindAll() ([]User, error)
	FindById(id int) (User, error)
	FindByEmail(email string) (User, error)
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

func (s *Store) FindAll() ([]User, error) {
	users := []User{}
	query := "SELECT * FROM users"
	err := s.db.Select(&users, query)

	return users, err
}

func (s *Store) FindById(id int) (User, error) {
	user := User{}
	query := "SELECT * FROM users WHERE id=$1"
	err := s.db.Get(&user, query, id)

	return user, err
}

func (s *Store) FindByEmail(email string) (User, error) {
	user := User{}
	query := "SELECT * FROM users WHERE email=$1"
	err := s.db.Get(&user, query, email)

	return user, err
}

func (s *Store) Create(input map[string]interface{}) error {
	columns := make([]string, 0, len(input))
	placeholders := make([]string, 0, len(input))
	values := make([]interface{}, 0, len(input))
	for col, val := range input {
		columns = append(columns, sqlx.Rebind(sqlx.DOLLAR, col))
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders)+1))
		values = append(values, val)
	}

	query := fmt.Sprintf("INSERT INTO users (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	_, err := s.db.Exec(query, values...)

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	columns := make([]string, 0, len(input))
	values := make([]interface{}, 0, len(input))
	for col, val := range input {
		columns = append(columns, fmt.Sprintf("%s=$%d", sqlx.Rebind(sqlx.DOLLAR, col), len(columns)+1))
		values = append(values, val)
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", strings.Join(columns, ", "), len(columns)+1)
	_, err := s.db.Exec(query, append(values, id)...)

	return err
}
