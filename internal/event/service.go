package event

import (
	"context"
	"time"

	"github.com/cinema-booker/third_party/tmdb"
)

type EventService interface {
	GetAll(ctx context.Context, pagination map[string]int) ([]Event, error)
	Get(ctx context.Context, id int) (Event, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
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

func (s *Service) GetAll(ctx context.Context, pagination map[string]int) ([]Event, error) {
	return s.store.FindAll(pagination)
}

func (s *Service) Get(ctx context.Context, id int) (Event, error) {
	return s.store.FindById(id)
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	return s.store.Create(input)
}

func (s *Service) Update(ctx context.Context, id int, input map[string]interface{}) error {
	return s.store.Update(id, input)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *Service) Restore(ctx context.Context, id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": nil,
	})
}
