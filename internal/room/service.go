package room

import "context"

type RoomService interface {
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}

type Service struct {
	store RoomStore
}

func NewService(store RoomStore) *Service {
	return &Service{
		store: store,
	}
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

	return s.store.Delete(id)
}
