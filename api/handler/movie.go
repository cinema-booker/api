package handler

import (
	"net/http"

	"github.com/cinema-booker/internal/movie"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/json"
	"github.com/gorilla/mux"
)

type MovieHandler struct {
	service movie.MovieService
}

func NewMovieHandler(service movie.MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

func (h *MovieHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/movies", errors.ErrorHandler(h.Search)).Methods(http.MethodGet)
}

func (h *MovieHandler) Search(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("query")

	event, err := h.service.Search(query)
	if err != nil {
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
