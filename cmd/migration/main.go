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
	var envConf = flag.String("conf", "config/local.yaml", "config path, eg: -conf ./config/local.yaml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := repository.NewDB(conf, log)
	migrateServer := migration.NewMigrateServer(db, log)

	defer func() {
		if err := migrateServer.Stop(context.Background()); err != nil {
			panic(err)
		}
	}()

	if err := migrateServer.Start(context.Background()); err != nil {
		panic(err)
	}
}
