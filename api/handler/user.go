package handler

import (
	"net/http"
	"strconv"

	"github.com/cinema-booker/api/middleware"
	"github.com/cinema-booker/api/utils"
	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/json"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service   user.UserService
	userStore user.UserStore
}

func NewUserHandler(service user.UserService, userStore user.UserStore) *UserHandler {
	return &UserHandler{
		service:   service,
		userStore: userStore,
	}
}

func (h *UserHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/users", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/users/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/users", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userStore))).Methods(http.MethodPost)
	mux.Handle("/users/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userStore))).Methods(http.MethodPatch)
	mux.Handle("/users/{id}", errors.ErrorHandler(middleware.IsAuth(h.Delete, h.userStore))).Methods(http.MethodDelete)
	mux.Handle("/users/{id}/restore", errors.ErrorHandler(middleware.IsAuth(h.Restore, h.userStore))).Methods(http.MethodPatch)
	mux.Handle("/users/{id}/password", errors.ErrorHandler(middleware.IsAuth(h.EditPassword, h.userStore))).Methods(http.MethodPatch)

	mux.Handle("/sign-up", errors.ErrorHandler(h.SignUp)).Methods(http.MethodPost)
	mux.Handle("/sign-in", errors.ErrorHandler(h.SignIn)).Methods(http.MethodPost)
	mux.Handle("/send-password-reset", errors.ErrorHandler(h.SendPasswordReset)).Methods(http.MethodPost)
	mux.Handle("/reset-password", errors.ErrorHandler(h.ResetPassword)).Methods(http.MethodPost)
	mux.Handle("/me", errors.ErrorHandler(middleware.IsAuth(h.GetMe, h.userStore))).Methods(http.MethodGet)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	pagination := utils.GetPaginationQueryParams(r)
	search := r.URL.Query().Get("search")

	users, err := h.service.GetAll(r.Context(), pagination, search)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, users); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, user); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Create(r.Context(), input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusCreated, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Update(r.Context(), id, input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusAccepted, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusNoContent, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) Restore(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Restore(r.Context(), id); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusNoContent, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.SignUp(r.Context(), input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusCreated, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	response, err := h.service.SignIn(r.Context(), input)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, response); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) SendPasswordReset(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	err := h.service.SendPasswordReset(r.Context(), input)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	err := h.service.ResetPassword(r.Context(), input)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) EditPassword(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	err = h.service.EditPassword(r.Context(), id, input)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) error {
	response, err := h.service.GetMe(r.Context())
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, response); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
