package tmdb

const (
	baseURL      = "https://api.themoviedb.org/3"
	imageBaseURL = "https://image.tmdb.org/t/p"
	imageSize    = "w500"
)

type TMDBMovie struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	OriginalLanguage string `json:"original_language"`
	PosterPath       string `json:"poster_path"`
}

type TMDBMovieDetails struct {
	Title            string `json:"title"`
	Overview         string `json:"overview"`
	OriginalLanguage string `json:"original_language"`
	PosterPath       string `json:"poster_path"`
	BackdropPath     string `json:"backdrop_path"`
}
