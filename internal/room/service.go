package room

import (
	"context"
	"time"
)

type RoomService interface {
	Create(ctx context.Context, cinemaId int, input map[string]interface{}) error
	Delete(ctx context.Context, cinemaId int, id int) error
}

type Service struct {
	store RoomStore
}

func NewService(store RoomStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(ctx context.Context, cinemaId int, input map[string]interface{}) error {
	input["cinema_id"] = cinemaId
	return s.store.Create(input)
}

func (s *Service) Delete(ctx context.Context, cinemaId int, id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
