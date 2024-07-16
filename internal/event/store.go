package event

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type EventStore interface {
	FindAll(pagination map[string]int, search string) ([]EventBasic, error)
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

func (s *Store) FindAll(pagination map[string]int, search string) ([]EventBasic, error) {
	events := []EventBasic{}

	offset := (pagination["page"] - 1) * pagination["limit"]
	query := `
    SELECT 
      e.id AS id,
			e.deleted_at AS deleted_at,
			c.id AS "cinema.id",
			c.user_id AS "cinema.user_id",
      c.name AS "cinema.name",
      c.description AS "cinema.description",
			c.deleted_at AS "cinema.deleted_at",
			a.id AS "cinema.address.id",
			a.address AS "cinema.address.address",
			a.longitude AS "cinema.address.longitude",
			a.latitude AS "cinema.address.latitude",
      m.id AS "movie.id",
      m.title AS "movie.title",
      m.description AS "movie.description",
      m.language AS "movie.language",
      m.poster AS "movie.poster",
      m.backdrop AS "movie.backdrop"
    FROM events e
    LEFT JOIN cinemas c ON e.cinema_id = c.id
    LEFT JOIN addresses a ON c.address_id = a.id
    LEFT JOIN movies m ON e.movie_id = m.id
		WHERE (
			c.name ILIKE '%' || $1 || '%'
			OR m.title ILIKE '%' || $1 || '%'
		)
		LIMIT $2 OFFSET $3
  `
	err := s.db.Select(&events, query, search, pagination["limit"], offset)

	return events, err
}

func (s *Store) FindById(id int) (Event, error) {
	event := Event{}
	query := `
    WITH booked_seats AS (
			SELECT
				s.id AS session_id,
				json_agg(b.place ORDER BY b.place) AS seats
			FROM sessions s
			LEFT JOIN	bookings b
				ON b.session_id = s.id AND b.status IN ('PENDING', 'CONFIRMED')
			GROUP BY s.id
    )
    SELECT
			e.id AS id,
			e.deleted_at AS deleted_at,
			c.id AS "cinema.id",
			c.user_id AS "cinema.user_id",
			c.name AS "cinema.name",
			c.description AS "cinema.description",
			c.deleted_at AS "cinema.deleted_at",
			a.id AS "cinema.address.id",
			a.address AS "cinema.address.address",
			a.longitude AS "cinema.address.longitude",
			a.latitude AS "cinema.address.latitude",
			m.id AS "movie.id",
			m.title AS "movie.title",
			m.description AS "movie.description",
			m.language AS "movie.language",
			m.poster AS "movie.poster",
			m.backdrop AS "movie.backdrop",
			COALESCE(
				json_agg(
					json_build_object(
						'id', s.id,
						'price', s.price,
						'starts_at', s.starts_at,
						'seats', COALESCE(bs.seats, '[]'),
						'room', json_build_object(
							'id', r.id,
							'number', r.number,
							'type', r.type
						)
					)
				) FILTER (WHERE s.id IS NOT NULL), '[]'
			) AS sessions
    FROM events e
    LEFT JOIN sessions s ON s.event_id = e.id
    LEFT JOIN booked_seats bs ON bs.session_id = s.id
    LEFT JOIN rooms r ON s.room_id = r.id
    LEFT JOIN cinemas c ON e.cinema_id = c.id
    LEFT JOIN addresses a ON c.address_id = a.id
    LEFT JOIN movies m ON e.movie_id = m.id
    WHERE e.id = $1
    GROUP BY e.id, c.id, a.id, m.id
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
		Title:       input["movie_title"].(string),
		Description: input["movie_description"].(string),
		Language:    input["movie_language"].(string),
		Poster:      input["movie_poster"].(string),
		Backdrop:    input["movie_backdrop"].(string),
	}

	var movieId int
	movieQuery := "INSERT INTO movies (title, description, language, poster, backdrop) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = tx.QueryRowx(movieQuery, movie.Title, movie.Description, movie.Language, movie.Poster, movie.Backdrop).Scan(&movieId)
	if err != nil {
		return err
	}

	query := "INSERT INTO events (cinema_id, movie_id) VALUES ($1, $2)"
	_, err = tx.Exec(query, input["cinema_id"], movieId)

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
