package utils

import (
	"net/http"
	"strconv"
)

func GetPaginationQueryParams(r *http.Request) map[string]int {
	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = 10
	}

	return map[string]int{
		"page":  page,
		"limit": limit,
	}
}
