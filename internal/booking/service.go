package booking

import (
	"time"
)

type BookingService interface {
	GetAll() ([]Booking, error)
	Get(id int) (Booking, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	Cancel(id int) error
}

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll() ([]Booking, error) {
	return s.store.FindAll()
}

func (s *Service) Get(id int) (Booking, error) {
	return s.store.FindById(id)
}

func (s *Service) Create(input map[string]interface{}) error {
	return s.store.Create(input)
}

func (s *Service) Update(id int, input map[string]interface{}) error {
	return s.store.Update(id, input)
}

func (s *Service) Cancel(id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"canceled_at": time.Now(),
	})
}
