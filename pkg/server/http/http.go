package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	*chi.Mux
	srv    *http.Server
	host   string
	port   int
	logger *slog.Logger
}

type Option func(s *Server)

func NewServer(mux *chi.Mux, logger *slog.Logger, args ...Option) *Server {
	s := &Server{
		Mux:    mux,
		logger: logger,
	}

	for _, opts := range args {
		opts(s)
	}
	return s
}

func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
func WithPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func (s *Server) Start(callerCtx context.Context) error {
	ctx, stop := signal.NotifyContext(callerCtx, os.Interrupt)
	defer stop()

	s.srv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s,
	}

	srvErr := make(chan error, 1)
	go func() {
		s.logger.Info("Starting server...")
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
	}()

	// We need to handle and return two different types of errors:
	// - received an error while shutting down the server
	// - received an error during server startup
	select {
	case <-ctx.Done():
		shutdownCtx, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdown()

		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			s.logger.Error("Error shutting down the HTTP server", "error", err)
			return err
		}
		return nil
	case err := <-srvErr:
		if err != nil {
			s.logger.Error("Error starting the HTTP server", "error", err)
			return err
		}
	}

	return nil
}
