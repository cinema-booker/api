package user

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cinema-booker/internal/constants"
	"github.com/cinema-booker/pkg/generator"
	"github.com/cinema-booker/pkg/hasher"
	"github.com/cinema-booker/pkg/jwt"
	"github.com/cinema-booker/third_party/resend"
)

type UserService interface {
	GetAll(ctx context.Context, pagination map[string]int) ([]User, error)
	Get(ctx context.Context, id int) (UserBasic, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	Restore(ctx context.Context, id int) error

	SignUp(ctx context.Context, input map[string]interface{}) error
	SignIn(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
	SendPasswordReset(ctx context.Context, input map[string]interface{}) error
	ResetPassword(ctx context.Context, input map[string]interface{}) error
	GetMe(ctx context.Context) (map[string]interface{}, error)
}

type Service struct {
	store UserStore
}

func NewService(store UserStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll(ctx context.Context, pagination map[string]int) ([]User, error) {
	return s.store.FindAll(pagination)
}

func (s *Service) Get(ctx context.Context, id int) (UserBasic, error) {
	return s.store.FindMeById(id)
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

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *Service) Restore(ctx context.Context, id int) error {
	_, err := s.store.FindById(id)
	if err != nil {
		return err
	}

	return s.store.Update(id, map[string]interface{}{
		"deleted_at": nil,
	})
}

func (s *Service) SignUp(ctx context.Context, input map[string]interface{}) error {
	_, err := s.store.FindByEmail(input["email"].(string))
	if err == nil {
		return fmt.Errorf("email already exists")
	}

	hashedPassword, err := hasher.Hash(input["password"].(string))
	if err != nil {
		return err
	}
	input["password"] = hashedPassword

	return s.store.Create(input)
}

func (s *Service) SignIn(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	user, err := s.store.FindByEmail(input["email"].(string))
	if err != nil {
		return nil, err
	}

	if !hasher.Compare(user.Password, input["password"].(string)) {
		return nil, fmt.Errorf("invalid credentials")
	}

	expiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		return nil, err
	}

	token, err := jwt.Create(os.Getenv("JWT_SECRET"), expiresIn, user.Id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"token": token,
	}, nil
}

func (s *Service) SendPasswordReset(ctx context.Context, input map[string]interface{}) error {
	user, err := s.store.FindByEmail(input["email"].(string))
	if err != nil {
		return err
	}

	generatedCode := generator.GenerateRandomCode(8)
	err = s.store.Update(user.Id, map[string]interface{}{
		"code":            generatedCode,
		"code_expires_at": time.Now().Add(time.Minute * 1),
	})
	if err != nil {
		return err
	}

	resend := resend.NewResendService(
		os.Getenv("RESEND_API_KEY"),
		os.Getenv("RESEND_FROM_EMAIL"),
	)

	_, err = resend.SendPasswordResetEmail([]string{user.Email}, generatedCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, input map[string]interface{}) error {
	user, err := s.store.FindByEmail(input["email"].(string))
	if err != nil {
		return err
	}

	if user.Code != input["code"].(string) {
		return fmt.Errorf("invalid code")
	}

	// TODO: `CodeExpiresAt` is not working well
	if user.CodeExpiresAt == nil || user.CodeExpiresAt.Before(time.Now()) {
		return fmt.Errorf("code expired")
	}

	hashedPassword, err := hasher.Hash(input["password"].(string))
	if err != nil {
		return err
	}

	err = s.store.Update(user.Id, map[string]interface{}{
		"password":        hashedPassword,
		"code":            "",
		"code_expires_at": nil,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetMe(ctx context.Context) (map[string]interface{}, error) {
	userId, ok := ctx.Value(constants.UserIDKey).(int)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}

	user, err := s.store.FindMeById(userId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":        user.Id,
		"name":      user.Name,
		"email":     user.Email,
		"role":      user.Role,
		"cinema_id": user.CinemaId,
	}, nil
}
