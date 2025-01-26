package main

import (
	"context"
	"flag"
	"github.com/giortzisg/go-boilerplate/internal/repository"
	"github.com/giortzisg/go-boilerplate/pkg/config"
	"github.com/giortzisg/go-boilerplate/pkg/migration"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var env = flag.String("config", "local", "config path, eg: -config local")
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
