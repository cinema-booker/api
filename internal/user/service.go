package user

import (
	"errors"
	"time"

	"github.com/cinema-booker/api/internal/auth/jwt"
	"github.com/cinema-booker/api/internal/auth/password"
)

type UserService interface {
	GetAll() ([]User, error)
	Get(id int) (User, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	Delete(id int) error
	SignUp(input map[string]interface{}) error
	SignIn(input map[string]interface{}) (string, error)
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
	user, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	if user.DeletedAt.Valid {
		return errors.New("user already deleted")
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *Service) SignUp(input map[string]interface{}) error {
	// TODO: validate input
	// required fields: first_name, last_name, email, password, role
	// role : MANAGER, VIEWER
	// email must be a valid email address
	// password must be at least 8 characters long
	// first_name and last_name must be at least 2 characters long

	hashedPassword, err := password.Hash(input["password"].(string))
	if err != nil {
		return err
	}
	input["password"] = hashedPassword

	return s.store.Create(input)
}

func (s *Service) SignIn(input map[string]interface{}) (string, error) {
	// TODO: validate input
	// required fields: email, password
	// email must be a valid email address
	// password must be at least 8 characters long

	user, err := s.store.FindByEmail(input["email"].(string))
	if err != nil {
		return "", err
	}

	if !password.Compare(user.Password, input["password"].(string)) {
		return "", errors.New("invalid password")
	}

	token, err := jwt.Create("secret", 3600*24*7, map[string]interface{}{
		"id": user.Id,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
