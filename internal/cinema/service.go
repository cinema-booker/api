package cinema

import (
	"fmt"
	"time"
)

type CinemaService interface {
	GetAll() ([]Cinema, error)
	Get(id int) (CinemaWithRooms, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	Delete(id int) error
}

type Service struct {
	store CinemaStore
}

func NewService(store CinemaStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll() ([]Cinema, error) {
	return s.store.FindAll()
}

func (s *Service) Get(id int) (CinemaWithRooms, error) {
	return s.store.FindById(id)
}

func (s *Service) Create(input map[string]interface{}) error {
	return s.store.Create(input)
}

func (s *Service) Update(id int, input map[string]interface{}) error {
	return s.store.Update(id, input)
}

func (s *Service) Delete(id int) error {
	cinema, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	if cinema.DeletedAt.Valid {
		return fmt.Errorf("cinema already deleted")
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
