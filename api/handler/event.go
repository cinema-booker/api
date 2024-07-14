package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/cinema-booker/api/middleware"
	"github.com/cinema-booker/api/utils"
	"github.com/cinema-booker/internal/event"
	"github.com/cinema-booker/internal/session"
	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/json"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	service        event.EventService
	sessionService session.SessionService
	userStore      user.UserStore
}

func NewEventHandler(service event.EventService, sessionService session.SessionService, userStore user.UserStore) *EventHandler {
	return &EventHandler{
		service:        service,
		sessionService: sessionService,
		userStore:      userStore,
	}
}

func (h *EventHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/events", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/events/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/events", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userStore))).Methods(http.MethodPost)
	mux.Handle("/events/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userStore))).Methods(http.MethodPatch)
	mux.Handle("/events/{id}", errors.ErrorHandler(middleware.IsAuth(h.Delete, h.userStore))).Methods(http.MethodDelete)
	mux.Handle("/events/{id}/restore", errors.ErrorHandler(middleware.IsAuth(h.Restore, h.userStore))).Methods(http.MethodPatch)

	mux.Handle("/events/{eventId}/sessions", errors.ErrorHandler(middleware.IsAuth(h.CreateSession, h.userStore))).Methods(http.MethodPost)
	mux.Handle("/events/{eventId}/sessions/{sessionId}", errors.ErrorHandler(middleware.IsAuth(h.DeleteSession, h.userStore))).Methods(http.MethodDelete)
}

func (h *EventHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	pagination := utils.GetPaginationQueryParams(r)

	events, err := h.service.GetAll(r.Context(), pagination)
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := json.Write(w, http.StatusOK, events); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	event, err := h.service.Get(r.Context(), id)
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

	if err := json.Write(w, http.StatusOK, event); err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	return nil
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
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

func (h *EventHandler) Restore(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.service.Restore(r.Context(), id); err != nil {
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

func (h *EventHandler) CreateSession(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	eventId, err := strconv.Atoi(vars["eventId"])
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

	if err := h.sessionService.Create(r.Context(), eventId, input); err != nil {
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

func (h *EventHandler) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	eventId, err := strconv.Atoi(vars["eventId"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}
	sessionId, err := strconv.Atoi(vars["sessionId"])
	if err != nil {
		return errors.HTTPError{
			Code: http.StatusInternalServerError,
			Err:  err,
		}
	}

	if err := h.sessionService.Delete(r.Context(), eventId, sessionId); err != nil {
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
