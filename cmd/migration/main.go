package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/giortzisg/go-boilerplate/internal/repository"
	"github.com/giortzisg/go-boilerplate/pkg/config"
	"github.com/giortzisg/go-boilerplate/pkg/migration"
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

	db := repository.NewDB(conf, logger)
	migrateServer := migration.NewMigrateServer(db, logger)

	defer func() {
		if err := migrateServer.Stop(context.Background()); err != nil {
			panic(err)
		}
	}()

	if err := migrateServer.Start(context.Background()); err != nil {
		panic(err)
	}
}
