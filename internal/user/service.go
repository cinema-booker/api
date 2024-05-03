package user

type UserService interface {
	GetAll() ([]User, error)
	Get(id int) (User, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	Delete(id int) error
}

type Service struct {
	store UserStore
}

func NewService(store UserStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll() ([]User, error) {
	return s.store.FindAll()
}

func (s *Service) Get(id int) (User, error) {
	return s.store.FindById(id)
}

func (s *Service) Create(input map[string]interface{}) error {
	return s.store.Create(input)
}

func (s *Service) Update(id int, input map[string]interface{}) error {
	return s.store.Update(id, input)
}

func (s *Service) Delete(id int) error {
	return s.store.Delete(id)
}
