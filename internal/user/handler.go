package user

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cinema-booker/api/internal/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	service UserService
}

func NewHandler(service UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/users", utils.ErrorHandler(h.GetAll)).Methods(http.MethodGet)
	mux.Handle("/users/{id}", utils.ErrorHandler(h.Get)).Methods(http.MethodGet)
	mux.Handle("/users", utils.ErrorHandler(h.Create)).Methods(http.MethodPost)
	mux.Handle("/users/{id}", utils.ErrorHandler(h.Update)).Methods(http.MethodPatch)
	mux.Handle("/users/{id}", utils.ErrorHandler(h.Delete)).Methods(http.MethodDelete)
	mux.Handle("/sign-up", utils.ErrorHandler(h.SignUp)).Methods(http.MethodPost)
	mux.Handle("/sign-in", utils.ErrorHandler(h.SignIn)).Methods(http.MethodPost)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.GetAll()
	if err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := utils.WriteJSON(w, http.StatusOK, users); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	user, err := h.service.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.HTTPError{
				Code: http.StatusNotFound,
				Err:  err,
			}
		}
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := utils.WriteJSON(w, http.StatusOK, user); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := utils.ParseJSON(r, &input); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Create(input); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := utils.WriteJSON(w, http.StatusCreated, nil); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	var input map[string]interface{}
	if err := utils.ParseJSON(r, &input); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Update(id, input); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := utils.WriteJSON(w, http.StatusAccepted, nil); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Delete(id); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := utils.WriteJSON(w, http.StatusNoContent, nil); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := utils.ParseJSON(r, &input); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.SignUp(input); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := utils.WriteJSON(w, http.StatusCreated, nil); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := utils.ParseJSON(r, &input); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	token, err := h.service.SignIn(input)
	if err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token}); err != nil {
		return utils.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}
