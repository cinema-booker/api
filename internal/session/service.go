package session

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"github.com/cinema-booker/pkg/errors"
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
	err := s.store.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, eventId int, id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	err = s.store.Update(id, map[string]interface{}{
		"deleted_at": time.Now(),
	})
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
