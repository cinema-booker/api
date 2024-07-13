package cinema

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CinemaStore interface {
	FindAll(pagination map[string]int) ([]Cinema, error)
	FindById(id int) (CinemaWithRooms, error)
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

func (s *Store) FindAll(pagination map[string]int) ([]Cinema, error) {
	cinemas := []Cinema{}

	offset := (pagination["page"] - 1) * pagination["limit"]
	query := `
    SELECT 
      c.id AS id,
			c.user_id AS user_id,
      c.name AS name,
      c.description AS description,
			c.deleted_at AS deleted_at,
      a.id AS "address.id",
      a.address AS "address.address",
      a.longitude AS "address.longitude",
      a.latitude AS "address.latitude"
    FROM cinemas c
    JOIN addresses a ON c.address_id = a.id
		WHERE c.deleted_at IS NULL
		LIMIT $1 OFFSET $2
  `
	err := s.db.Select(&cinemas, query, pagination["limit"], offset)

	return cinemas, err
}

func (s *Store) FindById(id int) (CinemaWithRooms, error) {
	cinema := CinemaWithRooms{}
	query := `
    SELECT
      c.id AS id,
      c.name AS name,
      c.description AS description,
			c.deleted_at AS deleted_at,
      a.id AS "address.id",
      a.address AS "address.address",
      a.longitude AS "address.longitude",
      a.latitude AS "address.latitude",
      COALESCE(
        json_agg(
          json_build_object(
            'id', r.id,
            'number', r.number,
            'type', r.type
          )
        ) FILTER (WHERE r.id IS NOT NULL),
        '[]'
      ) AS rooms
    FROM cinemas c
    LEFT JOIN addresses a ON c.address_id = a.id
    LEFT JOIN rooms r ON c.id = r.cinema_id
    WHERE c.id = $1 AND c.deleted_at IS NULL
    GROUP BY c.id, a.id
  `
	err := s.db.Get(&cinema, query, id)

	return cinema, err
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

	address := Address{
		Address:   input["address_address"].(string),
		Longitude: input["address_longitude"].(float64),
		Latitude:  input["address_latitude"].(float64),
	}

	var addressId int
	addressQuery := "INSERT INTO addresses (address, longitude, latitude) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRowx(addressQuery, address.Address, address.Longitude, address.Latitude).Scan(&addressId)
	if err != nil {
		return err
	}

	cinema := Cinema{
		Name:        input["name"].(string),
		Description: input["description"].(string),
	}

	var cinemaId int
	cinemaQuery := "INSERT INTO cinemas (user_id, address_id, name, description) VALUES ($1, $2, $3, $4) RETURNING id"
	err = tx.QueryRowx(cinemaQuery, input["user_id"], addressId, cinema.Name, cinema.Description).Scan(&cinemaId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	columns := make([]string, 0, len(input))
	values := make([]interface{}, 0, len(input))
	for col, val := range input {
		columns = append(columns, fmt.Sprintf("%s=$%d", sqlx.Rebind(sqlx.DOLLAR, col), len(columns)+1))
		values = append(values, val)
	}

	query := fmt.Sprintf("UPDATE cinemas SET %s WHERE id=$%d", strings.Join(columns, ", "), len(columns)+1)
	_, err := s.db.Exec(query, append(values, id)...)

	return err
}
