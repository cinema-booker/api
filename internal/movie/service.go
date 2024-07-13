package movie

import (
	"context"

	"github.com/cinema-booker/third_party/tmdb"
)

type MovieService interface {
	Search(ctx context.Context, query string) ([]map[string]interface{}, error)
}

type Service struct {
	tmdbService *tmdb.TMDB
}

func NewService(tmdbService *tmdb.TMDB) *Service {
	return &Service{
		tmdbService: tmdbService,
	}
}

func (s *Service) Search(ctx context.Context, query string) ([]map[string]interface{}, error) {
	movies, err := s.tmdbService.SearchMovies(query)
	if err != nil {
		return nil, err
	}

	return movies, nil
}
