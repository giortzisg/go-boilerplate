package http

import (
	"github.com/giortzisg/go-boilerplate/internal/handlers"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	*chi.Mux
	userHandler handlers.UserHandler
}

func NewRouter(userHandler handlers.UserHandler) *Router {
	router := &Router{
		Mux:         chi.NewRouter(),
		userHandler: userHandler,
	}

	router.RegisterUserRoutes()
	return router
}

func (r *Router) RegisterUserRoutes() {
	r.userHandler.Create()

	r.Route("/users", func(chi chi.Router) {
		chi.Post("/", r.userHandler.Create().ServeHTTP)
		chi.Get("/", r.userHandler.GetByEmail().ServeHTTP)
		chi.Put("/", r.userHandler.Update().ServeHTTP)
	})
}
