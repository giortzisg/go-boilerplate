package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/giortzisg/go-boilerplate/internal/app"
	"github.com/giortzisg/go-boilerplate/internal/handlers"
	routerHttp "github.com/giortzisg/go-boilerplate/internal/http"
	"github.com/giortzisg/go-boilerplate/internal/repository"
	"github.com/giortzisg/go-boilerplate/pkg/config"
	"github.com/giortzisg/go-boilerplate/pkg/server/http"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var env = flag.String("config", "config/local.yaml", "config path, eg: -config config/local.yaml")
	flag.Parse()
	conf, err := config.NewConfig(*env)
	if err != nil {
		logger.Error("error loading config", "error", err)
		os.Exit(1)
	}

	sqlDB := repository.NewDB(conf, logger)
	repo := repository.NewRepository(logger, sqlDB)
	userRepo := repository.NewUserRepository(repo)
	userService := app.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(handlers.NewHandler(logger), userService)
	router := routerHttp.NewRouter(logger, *userHandler)

	s := http.NewServer(
		router.Mux,
		logger,
		http.WithHost("localhost"),
		http.WithPort(8080),
	)

	if err := s.Start(context.Background()); err != nil {
		os.Exit(1)
	}
}
