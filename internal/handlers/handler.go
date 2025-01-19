package handlers

import (
	"log/slog"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(
	logger *slog.Logger,
) *Handler {
	return &Handler{
		logger,
	}
}
