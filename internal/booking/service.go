package booking

import (
	"context"
	"database/sql"
	goErrors "errors"

	"github.com/cinema-booker/internal/constants"
	"github.com/cinema-booker/internal/session"
	"github.com/cinema-booker/pkg/errors"
)

type BookingService interface {
	GetAll(ctx context.Context, pagination map[string]int, search string) ([]Booking, error)
	Get(ctx context.Context, id int) (Booking, error)
	Create(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
	Cancel(ctx context.Context, id int) error
	GetBookingWithUsersBySessionID(sessionID int, seats []string) (BookingWithUsers, error)
}

type Service struct {
	store        BookingStore
	sessionStore session.SessionStore
}

func NewService(store BookingStore, sessionStore session.SessionStore) *Service {
	return &Service{
		store:        store,
		sessionStore: sessionStore,
	}
}

func (s *Service) GetAll(ctx context.Context, pagination map[string]int, search string) ([]Booking, error) {
	bookings, err := s.store.FindAll(pagination, search)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return bookings, nil
}

func (s *Service) Get(ctx context.Context, id int) (Booking, error) {
	booking, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return booking, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return booking, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return booking, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	userId, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		return nil, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("user id not authenticated"),
		}
	}

	seatsInterface, ok := input["seats"].([]interface{})
	if !ok {
		return nil, errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("invalid seats input"),
		}
	}
	seats := make([]string, len(seatsInterface))
	for i, v := range seatsInterface {
		seats[i], ok = v.(string)
		if !ok {
			return nil, errors.CustomError{
				Key: errors.BadRequest,
				Err: goErrors.New("invalid seats input"),
			}
		}
	}

	sessionIdFloat, ok := input["session_id"].(float64)
	if !ok {
		return nil, errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("invalid session id input"),
		}
	}
	sessionId := int(sessionIdFloat)

	count, err := s.store.VerifySeatsCount(sessionId, seats)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	if count > 0 {
		return nil, errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("seats already booked"),
		}
	}

	for _, seat := range seats {
		err := s.store.Create(map[string]interface{}{
			"session_id": input["session_id"],
			"place":      seat,
			"user_id":    userId,
		})
		if err != nil {
			return nil, errors.CustomError{
				Key: errors.InternalServerError,
				Err: err,
			}
		}
	}

	session, err := s.sessionStore.FindById(sessionId)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return map[string]interface{}{
		"session_id": session.Id,
		"seats":      seats,
		"price":      int64(session.Price),
	}, nil
}

func (s *Service) Cancel(ctx context.Context, id int) error {
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
		"status": constants.BookingStatusCanceled,
	})
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) GetBookingWithUsersBySessionID(sessionID int, seats []string) (BookingWithUsers, error) {
	return s.store.ConfirmBookingBySessionAndSeats(sessionID, seats)
}
