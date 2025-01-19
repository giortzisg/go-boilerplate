package http

import (
	"github.com/spf13/viper"
	"log/slog"
)

type Server struct {
	logger *slog.Logger
	config *viper.Viper
}

func NewServer(logger *slog.Logger, config *viper.Viper) *Server {
	return &Server{
		logger,
		config,
	}
}
