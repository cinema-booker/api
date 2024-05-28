package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TMDB struct {
	APIKey string
}

func NewTMDBService(apiKey string) *TMDB {
	return &TMDB{
		APIKey: apiKey,
	}
}

func (t *TMDB) GetMovieById(id int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/movie/%d?api_key=%s&language=fr", baseURL, id, t.APIKey)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var movie TMDBMovieDetails
	if err = json.Unmarshal(body, &movie); err != nil {
		return nil, err
	}

	return formatMovieDetails(movie), nil
}

func (t *TMDB) SearchMovies(query string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s&page=1&language=fr", baseURL, t.APIKey, url.QueryEscape(query))

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Results []TMDBMovie `json:"results"`
	}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var movies []map[string]interface{}
	for _, movie := range result.Results {
		movies = append(movies, formatMovie(movie))
	}

	return movies, nil
}

func formatMovie(movie TMDBMovie) map[string]interface{} {
	return map[string]interface{}{
		"id":       movie.Id,
		"title":    movie.Title,
		"language": movie.OriginalLanguage,
		"poster":   fmt.Sprintf("%s/%s%s", imageBaseURL, imageSize, movie.PosterPath),
	}
}

func formatMovieDetails(movie TMDBMovieDetails) map[string]interface{} {
	return map[string]interface{}{
		"title":       movie.Title,
		"description": movie.Overview,
		"language":    movie.OriginalLanguage,
		"poster":      fmt.Sprintf("%s/%s%s", imageBaseURL, imageSize, movie.PosterPath),
		"backdrop":    fmt.Sprintf("%s/%s%s", imageBaseURL, imageSize, movie.BackdropPath),
	}
}
