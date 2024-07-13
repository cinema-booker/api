package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cinema-booker/api/middelware"
	"github.com/cinema-booker/internal/booking"
	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/json"
	"github.com/gorilla/mux"
)

type BookinHandler struct {
	service   booking.BookingService
	userStore user.UserStore
}

func NewBookingHandler(service booking.BookingService, userStore user.UserStore) *BookinHandler {
	return &BookinHandler{
		service:   service,
		userStore: userStore,
	}
}

func (h *BookinHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/bookings", errors.ErrorHandler(middelware.IsAuth(h.GetAll, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/bookings/{id}", errors.ErrorHandler(middelware.IsAuth(h.Get, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/bookings", errors.ErrorHandler(middelware.IsAuth(h.Create, h.userStore))).Methods(http.MethodPost)
	mux.Handle("/bookings/{id}", errors.ErrorHandler(middelware.IsAuth(h.Update, h.userStore))).Methods(http.MethodPatch)
	mux.Handle("/bookings/{id}", errors.ErrorHandler(middelware.IsAuth(h.Cancel, h.userStore))).Methods(http.MethodDelete)
}

func (h *BookinHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	bookings, err := h.service.GetAll(r.Context())
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, bookings); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *BookinHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	booking, err := h.service.Get(r.Context(), id)
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

	if err := json.Write(w, http.StatusOK, booking); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *BookinHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Create(r.Context(), input); err != nil {
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

func (h *BookinHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

	if err := h.service.Update(r.Context(), id, input); err != nil {
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

func (h *BookinHandler) Cancel(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Cancel(r.Context(), id); err != nil {
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
