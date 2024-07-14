package session

import (
	"context"
	"time"
)

type SessionService interface {
	Create(ctx context.Context, eventId int, input map[string]interface{}) error
	Delete(ctx context.Context, eventId int, id int) error
}

type Service struct {
	store SessionStore
}

func NewService(store SessionStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(ctx context.Context, eventId int, input map[string]interface{}) error {
	input["event_id"] = eventId
	return s.store.Create(input)
}

func (s *Service) Delete(ctx context.Context, eventId int, id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
