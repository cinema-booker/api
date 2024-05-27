package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cinema-booker/internal/cinema"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/json"
	"github.com/gorilla/mux"
)

type CinemaHandler struct {
	service cinema.CinemaService
}

func NewCinemaHandler(service cinema.CinemaService) *CinemaHandler {
	return &CinemaHandler{
		service: service,
	}
}

func (h *CinemaHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/cinemas", errors.ErrorHandler(h.GetAll)).Methods(http.MethodGet)
	mux.Handle("/cinemas/{id}", errors.ErrorHandler(h.Get)).Methods(http.MethodGet)
	mux.Handle("/cinemas", errors.ErrorHandler(h.Create)).Methods(http.MethodPost)
	mux.Handle("/cinemas/{id}", errors.ErrorHandler(h.Update)).Methods(http.MethodPatch)
	mux.Handle("/cinemas/{id}", errors.ErrorHandler(h.Delete)).Methods(http.MethodDelete)
	// TODO: Add routes for rooms
	// mux.Handle("/cinemas/{id}/rooms", errors.ErrorHandler(h.CreateRoom)).Methods(http.MethodPost)
	// mux.Handle("/cinemas/{id}/rooms/{roomId}", errors.ErrorHandler(h.UpdateRoom)).Methods(http.MethodPatch)
	// mux.Handle("/cinemas/{id}/rooms/{roomId}", errors.ErrorHandler(h.DeleteRoom)).Methods(http.MethodDelete)
}

func (h *CinemaHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	cinemas, err := h.service.GetAll()
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, cinemas); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *CinemaHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	cinema, err := h.service.Get(id)
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

	if err := json.Write(w, http.StatusOK, cinema); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *CinemaHandler) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *CinemaHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

func (h *CinemaHandler) Delete(w http.ResponseWriter, r *http.Request) error {
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
