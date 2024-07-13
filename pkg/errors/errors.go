package errors

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code int
	Err  error
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

type ErrorHandler func(w http.ResponseWriter, r *http.Request) error

func (f ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := f(w, r); err != nil {
		if e, ok := err.(HTTPError); ok {
			fmt.Println(e.Error())
			http.Error(w, e.Error(), e.Code)
		} else {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
