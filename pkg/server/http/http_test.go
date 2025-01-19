package http

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestServer_Start_Success(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := chi.NewRouter()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := NewServer(mux, logger, WithHost("127.0.0.1"), WithPort(8080))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(1 * time.Second)
		cancel() // Simulate server shutdown
	}()

	err := server.Start(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestServer_Shutdown_Gracefully(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := chi.NewRouter()
	server := NewServer(mux, logger, WithHost("127.0.0.1"), WithPort(8081))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(1 * time.Second)
		// Simulate receiving an interrupt signal
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(os.Interrupt)
	}()

	err := server.Start(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestServer_Start_ErrorOnStartup(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := chi.NewRouter()

	server := NewServer(mux, logger, WithHost("127.0.0.1"), WithPort(-1))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := server.Start(ctx)
	if err == nil {
		t.Error("expected error on startup, got nil")
	}
}

func TestServer_StartupPortAlreadyInUse(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := chi.NewRouter()
	server1 := NewServer(mux, logger, WithHost("127.0.0.1"), WithPort(8082))

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	go func() {
		_ = server1.Start(ctx1)
	}()

	time.Sleep(1 * time.Second)

	server2 := NewServer(mux, logger, WithHost("127.0.0.1"), WithPort(8082))

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	err := server2.Start(ctx2)
	if err == nil {
		t.Error("expected error for port already in use, got nil")
	}

	cancel1()
}

func TestServer_HandleOptions(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := chi.NewRouter()
	server := NewServer(mux, logger, WithHost("192.168.1.100"), WithPort(9090))

	if server.host != "192.168.1.100" {
		t.Errorf("expected host to be '192.168.1.100', got %v", server.host)
	}

	if server.port != 9090 {
		t.Errorf("expected port to be 9090, got %v", server.port)
	}
}
