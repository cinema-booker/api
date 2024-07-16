package room

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"github.com/cinema-booker/pkg/errors"
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
	err := s.store.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, cinemaId int, id int) error {
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
