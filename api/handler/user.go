package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/json"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service user.UserService
}

func NewUserHandler(service user.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/users", errors.ErrorHandler(h.GetAll)).Methods(http.MethodGet)
	mux.Handle("/users/{id}", errors.ErrorHandler(h.Get)).Methods(http.MethodGet)
	mux.Handle("/users", errors.ErrorHandler(h.Create)).Methods(http.MethodPost)
	mux.Handle("/users/{id}", errors.ErrorHandler(h.Update)).Methods(http.MethodPatch)
	mux.Handle("/users/{id}", errors.ErrorHandler(h.Delete)).Methods(http.MethodDelete)
	mux.Handle("/users/{id}/restore", errors.ErrorHandler(h.Restore)).Methods(http.MethodPatch)

	mux.Handle("/sign-up", errors.ErrorHandler(h.SignUp)).Methods(http.MethodPost)
	mux.Handle("/sign-in", errors.ErrorHandler(h.SignIn)).Methods(http.MethodPost)
	mux.Handle("/send-password-reset", errors.ErrorHandler(h.SendPasswordReset)).Methods(http.MethodPost)
	mux.Handle("/reset-password", errors.ErrorHandler(h.ResetPassword)).Methods(http.MethodPost)
	mux.Handle("/me", errors.ErrorHandler(h.GetMe)).Methods(http.MethodGet)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.GetAll()
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, users); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	user, err := h.service.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.HTTPError{
				Code: http.StatusNotFound,
				Err:  err,
			}
		}
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, user); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Create(input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusCreated, nil); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Update(id, input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusAccepted, nil); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Delete(id); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusNoContent, nil); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) Restore(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Restore(id); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusNoContent, nil); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.SignUp(input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusCreated, nil); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	response, err := h.service.SignIn(input)
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, response); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) SendPasswordReset(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	err := h.service.SendPasswordReset(input)
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, nil); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	err := h.service.ResetPassword(input)
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, nil); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token == "" {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("missing token"),
		}
	}

	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return errors.HTTPError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid token format"),
		}
	}

	response, err := h.service.GetMe(tokenParts[1])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, response); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}
