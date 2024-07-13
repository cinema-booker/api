package booking

import (
	"context"
	"time"
)

type BookingService interface {
	GetAll(ctx context.Context) ([]Booking, error)
	Get(ctx context.Context, id int) (Booking, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	Cancel(ctx context.Context, id int) error
}

type Service struct {
	store BookingStore
}

func NewService(store BookingStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]Booking, error) {
	return s.store.FindAll()
}

func (s *Service) Get(ctx context.Context, id int) (Booking, error) {
	return s.store.FindById(id)
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	return s.store.Create(input)
}

func (s *Service) Update(ctx context.Context, id int, input map[string]interface{}) error {
	return s.store.Update(id, input)
}

func (s *Service) Cancel(ctx context.Context, id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"canceled_at": time.Now(),
	})
}
