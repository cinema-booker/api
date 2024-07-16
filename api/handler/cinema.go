package handler

import (
	"net/http"
	"strconv"

	"github.com/cinema-booker/api/middleware"
	"github.com/cinema-booker/api/utils"
	"github.com/cinema-booker/internal/cinema"
	"github.com/cinema-booker/internal/room"
	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/json"
	"github.com/gorilla/mux"
)

type CinemaHandler struct {
	service     cinema.CinemaService
	roomService room.RoomService
	userStore   user.UserStore
}

func NewCinemaHandler(service cinema.CinemaService, roomService room.RoomService, userStore user.UserStore) *CinemaHandler {
	return &CinemaHandler{
		service:     service,
		roomService: roomService,
		userStore:   userStore,
	}
}

func (h *CinemaHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/cinemas", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/cinemas/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/cinemas", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userStore))).Methods(http.MethodPost)
	mux.Handle("/cinemas/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userStore))).Methods(http.MethodPatch)
	mux.Handle("/cinemas/{id}", errors.ErrorHandler(middleware.IsAuth(h.Delete, h.userStore))).Methods(http.MethodDelete)
	mux.Handle("/cinemas/{id}/restore", errors.ErrorHandler(middleware.IsAuth(h.Restore, h.userStore))).Methods(http.MethodPatch)

	mux.Handle("/cinemas/{cinemaId}/rooms", errors.ErrorHandler(middleware.IsAuth(h.CreateRoom, h.userStore))).Methods(http.MethodPost)
	mux.Handle("/cinemas/{cinemaId}/rooms/{roomId}", errors.ErrorHandler(middleware.IsAuth(h.DeleteRoom, h.userStore))).Methods(http.MethodDelete)
}

func (h *CinemaHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	pagination := utils.GetPaginationQueryParams(r)
	search := r.URL.Query().Get("search")

	cinemas, err := h.service.GetAll(r.Context(), pagination, search)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, cinemas); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *CinemaHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	cinema, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, cinema); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *CinemaHandler) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *CinemaHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

func (h *CinemaHandler) Delete(w http.ResponseWriter, r *http.Request) error {
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

func (h *CinemaHandler) Restore(w http.ResponseWriter, r *http.Request) error {
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

func (h *CinemaHandler) CreateRoom(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	cinemaId, err := strconv.Atoi(vars["cinemaId"])
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

	if err := h.roomService.Create(r.Context(), cinemaId, input); err != nil {
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

func (h *CinemaHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	cinemaId, err := strconv.Atoi(vars["cinemaId"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	roomId, err := strconv.Atoi(vars["roomId"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.roomService.Delete(r.Context(), cinemaId, roomId); err != nil {
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
