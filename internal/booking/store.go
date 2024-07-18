package booking

import (
	"fmt"
	"strings"

	"github.com/cinema-booker/internal/constants"
	"github.com/jmoiron/sqlx"
)

type BookingStore interface {
	FindAll(pagination map[string]int, search string) ([]Booking, error)
	FindById(id int) (Booking, error)
	VerifySeatsCount(sessionId int, seats []string) (int, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	ConfirmBookingBySessionAndSeats(sessionID int, seats []string) (BookingWithUsers, error)
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) FindAll(pagination map[string]int, search string) ([]Booking, error) {
	bookings := []Booking{}

	offset := (pagination["page"] - 1) * pagination["limit"]
	query := `
		SELECT 
      b.id AS id,
      b.place AS place,
			b.status AS status,
			u.id AS "user.id",
			u.name AS "user.name",
			s.id AS "session.id",
			s.price AS "session.price",
      s.starts_at AS "session.starts_at",
      r.id AS "session.room.id",
      r.number AS "session.room.number",
      r.type AS "session.room.type",
			e.id AS "session.event.id",
			c.id AS "session.event.cinema.id",
			c.name AS "session.event.cinema.name",
			c.description AS "session.event.cinema.description",
			c.user_id AS "session.event.cinema.user_id",
			c.deleted_at AS "session.event.cinema.deleted_at",
			a.id AS "session.event.cinema.address.id",
			a.address AS "session.event.cinema.address.address",
			a.longitude AS "session.event.cinema.address.longitude",
			a.latitude AS "session.event.cinema.address.latitude",
			m.id AS "session.event.movie.id",
			m.title AS "session.event.movie.title",
			m.description AS "session.event.movie.description",
			m.language AS "session.event.movie.language",
			m.poster AS "session.event.movie.poster",
			m.backdrop AS "session.event.movie.backdrop"
    FROM bookings b
    LEFT JOIN users u ON b.user_id = u.id
    LEFT JOIN sessions s ON b.session_id = s.id
    LEFT JOIN rooms r ON s.room_id = r.id
    LEFT JOIN events e ON s.event_id = e.id
    LEFT JOIN cinemas c ON e.cinema_id = c.id
    LEFT JOIN movies m ON e.movie_id = m.id
    LEFT JOIN addresses a ON c.address_id = a.id
		WHERE (
			u.name ILIKE '%' || $1 || '%'
			OR c.name ILIKE '%' || $1 || '%'
			OR m.title ILIKE '%' || $1 || '%'
		)
		LIMIT $2 OFFSET $3
	`
	err := s.db.Select(&bookings, query, search, pagination["limit"], offset)

	return bookings, err
}

func (s *Store) FindById(id int) (Booking, error) {
	booking := Booking{}
	query := `
		SELECT 
      b.id AS id,
      b.place AS place,
			b.status AS status,
			u.id AS "user.id",
			u.name AS "user.name",
			s.id AS "session.id",
			s.price AS "session.price",
      s.starts_at AS "session.starts_at",
      r.id AS "session.room.id",
      r.number AS "session.room.number",
      r.type AS "session.room.type",
			e.id AS "session.event.id",
			c.id AS "session.event.cinema.id",
			c.name AS "session.event.cinema.name",
			c.description AS "session.event.cinema.description",
			c.user_id AS "session.event.cinema.user_id",
			c.deleted_at AS "session.event.cinema.deleted_at",
			a.id AS "session.event.cinema.address.id",
			a.address AS "session.event.cinema.address.address",
			a.longitude AS "session.event.cinema.address.longitude",
			a.latitude AS "session.event.cinema.address.latitude",
			m.id AS "session.event.movie.id",
			m.title AS "session.event.movie.title",
			m.description AS "session.event.movie.description",
			m.language AS "session.event.movie.language",
			m.poster AS "session.event.movie.poster",
			m.backdrop AS "session.event.movie.backdrop"
    FROM bookings b
    LEFT JOIN users u ON b.user_id = u.id
    LEFT JOIN sessions s ON b.session_id = s.id
    LEFT JOIN rooms r ON s.room_id = r.id
    LEFT JOIN events e ON s.event_id = e.id
    LEFT JOIN cinemas c ON e.cinema_id = c.id
    LEFT JOIN movies m ON e.movie_id = m.id
    LEFT JOIN addresses a ON c.address_id = a.id
		WHERE b.id=$1
	`
	err := s.db.Get(&booking, query, id)

	return booking, err
}

func (s *Store) VerifySeatsCount(sessionId int, seats []string) (int, error) {
	query := `
		SELECT COUNT(DISTINCT b.id)
		FROM bookings b
		WHERE b.place IN (?) AND b.session_id = ? AND b.status IN (?)
  `
	query, args, err := sqlx.In(query, seats, sessionId, []string{
		constants.BookingStatusPending,
		constants.BookingStatusConfirmed,
	})
	if err != nil {
		return -1, err
	}

	query = s.db.Rebind(query)
	var count int
	err = s.db.Get(&count, query, args...)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (s *Store) Create(input map[string]interface{}) error {
	query := "INSERT INTO bookings (user_id, session_id, place) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(query, input["user_id"], input["session_id"], input["place"])

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

func (s *Store) ConfirmBookingBySessionAndSeats(sessionID int, seats []string) (BookingWithUsers, error) {

	queryUpdate := `
	UPDATE bookings 
	SET status = constants.BookingStatusConfirmed 
	WHERE session_id = ? AND place IN (?)
	`
	queryUpdate, argsUpdate, err := sqlx.In(queryUpdate, sessionID, seats)
	if err != nil {
		return BookingWithUsers{}, err
	}
	queryUpdate = s.db.Rebind(queryUpdate)

	_, err = s.db.Exec(queryUpdate, argsUpdate...)
	if err != nil {
		return BookingWithUsers{}, err
	}

	querySelect := `
	SELECT 
		b.id as booking_id, b.place, b.status, 
		u.id as booking_user_id, u.name as booking_user_name,
		cu.id as cinema_user_id, cu.name as cinema_user_name, cu.email as cinema_user_email, cu.role as cinema_user_role 
	FROM 
		bookings b
	JOIN 
		users u ON b.user_id = u.id
	JOIN 
		sessions s ON b.session_id = s.id
	JOIN 
		rooms r ON s.room_id = r.id
	JOIN 
		cinemas c ON r.cinema_id = c.id
	JOIN 
		users cu ON c.user_id = cu.id
	WHERE 
		b.session_id = ? AND b.place IN (?)
	LIMIT 1
	`
	querySelect, argsSelect, err := sqlx.In(querySelect, sessionID, seats)
	if err != nil {
		return BookingWithUsers{}, err
	}
	querySelect = s.db.Rebind(querySelect)

	var result BookingWithUsers

	err = s.db.QueryRow(querySelect, argsSelect...).Scan(
		&result.Booking.Id, &result.Booking.Place, &result.Booking.Status,
		&result.BookingUser.Id, &result.BookingUser.Name,
		&result.CinemaUser.Id, &result.CinemaUser.Name, &result.CinemaUser.Email, &result.CinemaUser.Role,
	)
	if err != nil {
		return result, err
	}

	return result, nil
}
