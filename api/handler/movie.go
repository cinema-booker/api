package handler

import (
	"net/http"

	"github.com/cinema-booker/api/middleware"
	"github.com/cinema-booker/internal/movie"
	"github.com/cinema-booker/internal/user"
	"github.com/cinema-booker/pkg/errors"
	"github.com/cinema-booker/pkg/json"
	"github.com/gorilla/mux"
)

type MovieHandler struct {
	service   movie.MovieService
	userStore user.UserStore
}

func NewMovieHandler(service movie.MovieService, userStore user.UserStore) *MovieHandler {
	return &MovieHandler{
		service:   service,
		userStore: userStore,
	}
}

func (h *MovieHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/movies", errors.ErrorHandler(middleware.IsAuth(h.Search, h.userStore))).Methods(http.MethodGet)
}

func (h *MovieHandler) Search(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("query")

	event, err := h.service.Search(r.Context(), query)
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
