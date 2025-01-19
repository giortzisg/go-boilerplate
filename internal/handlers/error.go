package handlers

import (
	"errors"
	v1 "github.com/giortzisg/go-boilerplate/api/v1"
	"github.com/giortzisg/go-boilerplate/pkg/json"
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

			response := &v1.Response{
				Code:    status,
				Message: err.Error(),
				Data:    nil,
			}

			if err = json.Encoder(w, response, status); err != nil {
				// if encoding fails, fallback to writing the error message directly
				http.Error(w, err.Error(), status)
			}
		}
	})
}
