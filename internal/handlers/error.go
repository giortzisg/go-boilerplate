package handlers

import (
	"errors"
	"net/http"
)

func ErrorHandler(f func(http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			var status int
			var statusErr interface {
				error
				HTTPStatus() int
			}
			if errors.As(err, &statusErr) {
				status = statusErr.HTTPStatus()
			}
			http.Error(w, err.Error(), status)
		}
	})
}
