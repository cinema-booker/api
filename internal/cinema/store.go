package cinema

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CinemaStore interface {
	FindAll() ([]Cinema, error)
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

func (s *Store) FindAll() ([]Cinema, error) {
	cinemas := []Cinema{}
	query := `
    SELECT 
      c.id AS id,
      c.name AS name,
      c.description AS description,
			c.deleted_at AS deleted_at,
      a.id AS "address.id",
      a.country AS "address.country",
      a.city AS "address.city",
      a.zip_code AS "address.zip_code",
      a.street AS "address.street",
      a.longitude AS "address.longitude",
      a.latitude AS "address.latitude"
    FROM cinemas c
    JOIN addresses a ON c.address_id = a.id
  `
	err := s.db.Select(&cinemas, query)

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
      a.country AS "address.country",
      a.city AS "address.city",
      a.zip_code AS "address.zip_code",
      a.street AS "address.street",
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
    WHERE c.id = $1
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
		Country:   input["address_country"].(string),
		City:      input["address_city"].(string),
		ZipCode:   input["address_zip_code"].(string),
		Street:    input["address_street"].(string),
		Longitude: input["address_longitude"].(float64),
		Latitude:  input["address_latitude"].(float64),
	}

	var addressId int
	addressQuery := "INSERT INTO addresses (country, city, zip_code, street, longitude, latitude) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = tx.QueryRowx(addressQuery, address.Country, address.City, address.ZipCode, address.Street, address.Longitude, address.Latitude).Scan(&addressId)
	if err != nil {
		return err
	}

	cinema := Cinema{
		Name:        input["name"].(string),
		Description: input["description"].(string),
	}

	var cinemaId int
	cinemaQuery := "INSERT INTO cinemas (address_id, name, description) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRowx(cinemaQuery, addressId, cinema.Name, cinema.Description).Scan(&cinemaId)
	if err != nil {
		return err
	}

	userIdFloat64, ok := input["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid type for user_id")
	}
	userId := int(userIdFloat64)
	userCinemaQuery := "INSERT INTO users_cinemas (user_id, cinema_id) VALUES ($1, $2)"
	_, err = tx.Exec(userCinemaQuery, userId, cinemaId)

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
