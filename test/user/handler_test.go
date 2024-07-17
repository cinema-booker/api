package handler

import (
	"context"
	"os"

	// "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cinema-booker/api/handler"
	"github.com/cinema-booker/pkg/jwt"

	"github.com/cinema-booker/internal/constants"
	"github.com/cinema-booker/internal/user"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserService struct {
	mock.Mock
}

// Create implements user.UserService.
func (m *MockUserService) Create(ctx context.Context, input map[string]interface{}) error {
	return m.Called(ctx, input).Error(0)
}

// Delete implements user.UserService.
func (m *MockUserService) Delete(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

// Get implements user.UserService.
func (m *MockUserService) Get(ctx context.Context, id int) (user.UserBasic, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user.UserBasic), args.Error(1)
}

// GetAll implements user.UserService.
func (m *MockUserService) GetAll(ctx context.Context, pagination map[string]int, search string) ([]user.User, error) {
	args := m.Called(ctx, pagination)
	return args.Get(0).([]user.User), args.Error(1)
}

// GetMe implements user.UserService.
func (m *MockUserService) GetMe(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// ResetPassword implements user.UserService.
func (m *MockUserService) ResetPassword(ctx context.Context, input map[string]interface{}) error {
	return m.Called(ctx, input).Error(0)
}

// Restore implements user.UserService.
func (m *MockUserService) Restore(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

// SendPasswordReset implements user.UserService.
func (m *MockUserService) SendPasswordReset(ctx context.Context, input map[string]interface{}) error {
	return m.Called(ctx, input).Error(0)
}

// SignIn implements user.UserService.
func (m *MockUserService) SignIn(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// SignUp implements user.UserService.
func (m *MockUserService) SignUp(ctx context.Context, input map[string]interface{}) error {
	return m.Called(ctx, input).Error(0)
}

// Update implements user.UserService.
func (m *MockUserService) Update(ctx context.Context, id int, input map[string]interface{}) error {
	return m.Called(ctx, id, input).Error(0)
}

type MockUserStore struct {
	mock.Mock
}

// Create implements user.UserStore.
func (m *MockUserStore) Create(input map[string]interface{}) error {
	return m.Called(input).Error(0)
}

// FindAll implements user.UserStore.
func (m *MockUserStore) FindAll(pagination map[string]int, search string) ([]user.User, error) {
	args := m.Called(pagination)
	return args.Get(0).([]user.User), args.Error(1)
}

// FindByEmail implements user.UserStore.
func (m *MockUserStore) FindByEmail(email string) (user.User, error) {
	args := m.Called(email)
	return args.Get(0).(user.User), args.Error(1)
}

// FindById implements user.UserStore.
func (m *MockUserStore) FindById(id int) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}

// FindMeById implements user.UserStore.
func (m *MockUserStore) FindMeById(id int) (user.UserBasic, error) {
	args := m.Called(id)
	return args.Get(0).(user.UserBasic), args.Error(1)
}

// Update implements user.UserStore.
func (m *MockUserStore) Update(id int, input map[string]interface{}) error {
	return m.Called(id, input).Error(0)
}

// TestGetAll
func TestGetAll(t *testing.T) {
	mockService := new(MockUserService)
	mockStore := new(MockUserStore)
	userHandler := handler.NewUserHandler(mockService, mockStore)

	mockService.On("GetAll", mock.Anything, mock.Anything).Return([]user.User{}, nil)
	mockStore.On("FindById", mock.Anything).Return(user.User{Id: 1, Role: constants.UserRoleAdmin}, nil)

	// Mock token and context
	token, err := jwt.Create(os.Getenv("JWT_SECRET"), 3600, 1)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mock user ID and role
			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserIDKey, 1)
			ctx = context.WithValue(ctx, constants.UserRoleKey, constants.UserRoleAdmin)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	})
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

// TestGetUser
func TestGetUser(t *testing.T) {
	mockService := new(MockUserService)
	mockStore := new(MockUserStore)
	userHandler := handler.NewUserHandler(mockService, mockStore)

	mockService.On("Get", mock.Anything, 1).Return(user.UserBasic{Id: 1, Name: "Test User"}, nil)

	mockStore.On("FindById", 1).Return(user.User{Id: 1, Role: constants.UserRoleAdmin}, nil)

	req, err := http.NewRequest(http.MethodGet, "/users/1", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	token, err := jwt.Create(os.Getenv("JWT_SECRET"), 3600, 1)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserIDKey, 1)
			ctx = context.WithValue(ctx, constants.UserRoleKey, constants.UserRoleAdmin)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	})
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

// TestGetMe
func TestGetMe(t *testing.T) {
	mockService := new(MockUserService)
	mockStore := new(MockUserStore)
	userHandler := handler.NewUserHandler(mockService, mockStore)

	expectedResponse := map[string]interface{}{
		"id":        1,
		"name":      "John Doe",
		"email":     "johndoe@example.com",
		"role":      constants.UserRoleAdmin,
		"cinema_id": 1,
	}

	mockService.On("GetMe", mock.Anything).Return(expectedResponse, nil)
	mockStore.On("FindById", 1).Return(user.User{Id: 1, Role: constants.UserRoleAdmin}, nil)

	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	require.NoError(t, err)

	token, err := jwt.Create(os.Getenv("JWT_SECRET"), 3600, 1)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, constants.UserIDKey, 1)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	})
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}
