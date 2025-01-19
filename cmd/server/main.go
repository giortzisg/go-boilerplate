package main

import (
	"context"
	"flag"
	"github.com/giortzisg/go-boilerplate/internal/app"
	"github.com/giortzisg/go-boilerplate/internal/handlers"
	routerHttp "github.com/giortzisg/go-boilerplate/internal/http"
	"github.com/giortzisg/go-boilerplate/internal/repository"
	"github.com/giortzisg/go-boilerplate/pkg/config"
	"github.com/giortzisg/go-boilerplate/pkg/server/http"
	"log/slog"
	"os"
)

func main() {
	var envConf = flag.String("conf", "config/local.yaml", "config path, eg: -conf ./config/local.yaml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sqlDB := repository.NewDB(conf, logger)
	repo := repository.NewRepository(logger, sqlDB)
	userRepo := repository.NewUserRepository(repo)
	userService := app.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(handlers.NewHandler(logger), userService)
	router := routerHttp.NewRouter(*userHandler)

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
