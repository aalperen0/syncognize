package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/aalperen0/syncognize/internal/config"
	"go.uber.org/zap"
)

func NewPostgres(ctx context.Context, cfg config.DatabaseConfig, logger *zap.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	// fmt.Println("DSN", dsn)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("failed to open postgres connection", zap.Error(err))
		return nil, err
	}

	// Apply connection pool settings
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	if err := db.PingContext(ctx); err != nil {
		logger.Error("failed to ping postgres", zap.Error(err))
		return nil, err
	}

	logger.Info("postgres connected")

	return db, nil
}
