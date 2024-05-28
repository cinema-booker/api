package event

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/cinema-booker/third_party/tmdb"
)

type EventService interface {
	GetAll() ([]Event, error)
	Get(id int) (Event, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	Delete(id int) error
	Restore(id int) error
}

type Service struct {
	store       EventStore
	tmdbService *tmdb.TMDB
}

func NewService(store EventStore, tmdbService *tmdb.TMDB) *Service {
	return &Service{
		store:       store,
		tmdbService: tmdbService,
	}
}

func (s *Service) GetAll() ([]Event, error) {
	return s.store.FindAll()
}

func (s *Service) Get(id int) (Event, error) {
	return s.store.FindById(id)
}

func (s *Service) Create(input map[string]interface{}) error {
	movieIdFloat64, ok := input["movie_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid type for movie_id")
	}
	movieId := int(movieIdFloat64)

	movie, err := s.tmdbService.GetMovieById(movieId)
	if err != nil {
		return err
	}
	input["movie"] = movie

	return s.store.Create(input)
}

func (s *Service) Update(id int, input map[string]interface{}) error {
	return s.store.Update(id, input)
}

func (s *Service) Delete(id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (s *Service) Restore(id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	})
}
