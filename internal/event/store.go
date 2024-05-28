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
	query := `
    SELECT 
      e.id AS id,
      e.price AS price,
      e.starts_at AS starts_at,
			e.ends_at AS ends_at,
			e.deleted_at AS deleted_at,
      r.id AS "room.id",
      r.number AS "room.number",
      r.type AS "room.type",
      m.id AS "movie.id",
      m.title AS "movie.title",
      m.description AS "movie.description",
      m.language AS "movie.language",
      m.poster AS "movie.poster",
      m.backdrop AS "movie.backdrop"
    FROM events e
    LEFT JOIN rooms r ON e.room_id = r.id
    LEFT JOIN movies m ON e.movie_id = m.id
  `
	err := s.db.Select(&events, query)

	return events, err
}

func (s *Store) FindById(id int) (Event, error) {
	event := Event{}
	query := `
    SELECT 
      e.id AS id,
      e.price AS price,
      e.starts_at AS starts_at,
			e.ends_at AS ends_at,
			e.deleted_at AS deleted_at,
      r.id AS "room.id",
      r.number AS "room.number",
      r.type AS "room.type",
      m.id AS "movie.id",
      m.title AS "movie.title",
      m.description AS "movie.description",
      m.language AS "movie.language",
      m.poster AS "movie.poster",
      m.backdrop AS "movie.backdrop"
    FROM events e
    LEFT JOIN rooms r ON e.room_id = r.id
    LEFT JOIN movies m ON e.movie_id = m.id
		WHERE e.id=$1
  `
	err := s.db.Get(&event, query, id)

	return event, err
}

func (s *Store) Create(input map[string]interface{}) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	movie := Movie{
		Title:       input["movie"].(map[string]interface{})["title"].(string),
		Description: input["movie"].(map[string]interface{})["overview"].(string),
		Language:    input["movie"].(map[string]interface{})["language"].(string),
		Poster:      input["movie"].(map[string]interface{})["poster"].(string),
		Backdrop:    input["movie"].(map[string]interface{})["backdrop"].(string),
	}

	var movieId int
	movieQuery := "INSERT INTO movies (title, description, language, poster, backdrop) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = tx.QueryRowx(movieQuery, movie.Title, movie.Description, movie.Language, movie.Poster, movie.Backdrop).Scan(&movieId)
	if err != nil {
		return err
	}

	query := "INSERT INTO events (movie_id, room_id, price, starts_at, ends_at) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(query, movieId, input["room_id"], input["price"], input["starts_at"], input["ends_at"])

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
