package migrations

import (
	"context"
	"github.com/giortzisg/go-boilerplate/internal/model"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type MigrateServer struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewMigrateServer(db *gorm.DB, log *slog.Logger) *MigrateServer {
	return &MigrateServer{
		db:  db,
		log: log,
	}
}
func (m *MigrateServer) Start(ctx context.Context) error {
	if err := m.db.AutoMigrate(
		&model.User{},
	); err != nil {
		m.log.Warn("user migrate error", "err", err)
		return err
	}
	m.log.Info("AutoMigrate success")
	os.Exit(0)
	return nil
}
func (m *MigrateServer) Stop(ctx context.Context) error {
	m.log.Info("AutoMigrate stop")
	return nil
}
