package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/cinema-booker/internal/constants"
)

type BookingService interface {
	GetAll(ctx context.Context, pagination map[string]int) ([]Booking, error)
	Get(ctx context.Context, id int) (Booking, error)
	Create(ctx context.Context, input map[string]interface{}) error
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

func (s *Service) GetAll(ctx context.Context, pagination map[string]int) ([]Booking, error) {
	return s.store.FindAll(pagination)
}

func (s *Service) Get(ctx context.Context, id int) (Booking, error) {
	return s.store.FindById(id)
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	userId, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	seatsInterface, ok := input["seats"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid seats input")
	}
	seats := make([]string, len(seatsInterface))
	for i, v := range seatsInterface {
		seats[i], ok = v.(string)
		if !ok {
			return fmt.Errorf("invalid seat value at index %d", i)
		}
	}

	sessionIdFloat, ok := input["session_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid session_id input")
	}
	sessionId := int(sessionIdFloat)

	count, err := s.store.VerifySeatsCount(sessionId, seats)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("seats already booked")
	}

	for _, seat := range seats {
		err := s.store.Create(map[string]interface{}{
			"session_id": input["session_id"],
			"place":      seat,
			"user_id":    userId,
		})
		if err != nil {
			return err
		}
	}

	return nil

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
