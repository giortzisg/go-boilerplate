package handlers

import (
	v1 "github.com/giortzisg/go-boilerplate/api/v1"
	"github.com/giortzisg/go-boilerplate/internal/app"
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
	return ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		requestData, err := json.Decoder[v1.CreateUserRequest](r)
		if err != nil {
			return err
		}

		if err = h.userService.Create(r.Context(), requestData); err != nil {
			return err
		}

		return json.Encoder(
			w, &v1.Response{
				Message: "User created successfully",
				Code:    http.StatusCreated,
			},
			http.StatusCreated,
		)
	})
}

func (h *UserHandler) GetByEmail() http.Handler {
	return ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		requestData, err := json.Decoder[v1.GetUserByEmailRequest](r)
		if err != nil {
			return err
		}

		response, err := h.userService.GetByEmail(r.Context(), requestData)
		if err != nil {
			return err
		}

		return json.Encoder(
			w,
			&v1.Response{
				Message: "User retrieved successfully",
				Code:    http.StatusOK,
				Data:    response,
			},
			http.StatusOK,
		)
	})
}

func (h *UserHandler) Update() http.Handler {
	return ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		requestData, err := json.Decoder[v1.UpdateUserRequest](r)
		if err != nil {
			return err
		}

		if err = h.userService.Update(r.Context(), requestData); err != nil {
			return err
		}

		return json.Encoder(
			w,
			&v1.Response{
				Message: "User updated successfully",
				Code:    http.StatusOK,
			},
			http.StatusOK,
		)
	})
}
