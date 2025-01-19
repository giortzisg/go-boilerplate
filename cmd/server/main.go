package main

import (
	"context"
	"github.com/giortzisg/go-boilerplate/pkg/server/http"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	m := chi.NewRouter()

	s := http.NewServer(
		m,
		logger,
		http.WithHost("localhost"),
		http.WithPort(8080),
	)

	if err := s.Start(context.Background()); err != nil {
		os.Exit(1)
	}
}
