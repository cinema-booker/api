package cinema

import (
	"context"
	"fmt"
	"time"

	"github.com/cinema-booker/internal/constants"
)

type CinemaService interface {
	GetAll(ctx context.Context, pagination map[string]int) ([]Cinema, error)
	Get(ctx context.Context, id int) (CinemaWithRooms, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error
}

type Service struct {
	store CinemaStore
}

func NewService(store CinemaStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll(ctx context.Context, pagination map[string]int) ([]Cinema, error) {
	return s.store.FindAll(pagination)
}

func (s *Service) Get(ctx context.Context, id int) (CinemaWithRooms, error) {
	return s.store.FindById(id)
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	userId, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		return fmt.Errorf("Unauthorized")
	}
	input["user_id"] = userId

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
