package handlers

import (
	v1 "github.com/giortzisg/go-boilerplate/api/v1"
	"github.com/giortzisg/go-boilerplate/internal/app"
	e "github.com/giortzisg/go-boilerplate/pkg/error"
	"github.com/giortzisg/go-boilerplate/pkg/json"
	"net/http"
)

type UserHandler struct {
	*Handler
	userService app.UserService
}

func NewUserHandler(h *Handler, userService app.UserService) *UserHandler {
	return &UserHandler{
		Handler:     h,
		userService: userService,
	}
}

func (h *UserHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestData, err := json.Decoder[v1.CreateUserRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = h.userService.Create(r.Context(), requestData); err != nil {
			http.Error(w, err.Error(), err.(*e.StatusError).HTTPStatus())
			return
		}

		if err = json.Encoder(
			w, &v1.Response{
				Message: "User created successfully",
				Code:    http.StatusCreated,
			},
			http.StatusCreated,
		); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (h *UserHandler) GetByEmail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestData, err := json.Decoder[v1.GetUserByEmailRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := h.userService.GetByEmail(r.Context(), requestData)
		if err != nil {
			http.Error(w, err.Error(), err.(*e.StatusError).HTTPStatus())
			return
		}

		if err = json.Encoder(
			w,
			&v1.Response{
				Message: "User retrieved successfully",
				Code:    http.StatusOK,
				Data:    response,
			},
			http.StatusOK,
		); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (h *UserHandler) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestData, err := json.Decoder[v1.UpdateUserRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = h.userService.Update(r.Context(), requestData); err != nil {
			http.Error(w, err.Error(), err.(*e.StatusError).HTTPStatus())
			return
		}

		if err = json.Encoder(
			w,
			&v1.Response{
				Message: "User updated successfully",
				Code:    http.StatusOK,
			},
			http.StatusOK,
		); err != nil {
			http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
