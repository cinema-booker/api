package room

type RoomService interface {
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	Delete(id int) error
}

type Service struct {
	store RoomStore
}

func NewService(store RoomStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(input map[string]interface{}) error {
	return s.store.Create(input)
}

func (s *Service) Update(id int, input map[string]interface{}) error {
	return s.store.Update(id, input)
}

func (s *Service) Delete(id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Delete(id)
}
