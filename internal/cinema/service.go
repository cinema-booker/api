package cinema

import (
	"context"
	"database/sql"
	goErrors "errors"
	"time"

	"github.com/cinema-booker/internal/constants"
	"github.com/cinema-booker/pkg/errors"
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
	cinemas, err := s.store.FindAll(pagination)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return cinemas, nil
}

func (s *Service) Get(ctx context.Context, id int) (CinemaWithRooms, error) {
	cinema, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return cinema, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return cinema, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return cinema, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	userId, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not found in context"),
		}
	}

	input["user_id"] = userId
	err := s.store.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id int, input map[string]interface{}) error {
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

	err = s.store.Update(id, input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
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

func (s *Service) Restore(ctx context.Context, id int) error {
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
		"deleted_at": nil,
	})
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
