package main

import (
	"context"
	"fmt"

	"github.com/aalperen0/syncognize/internal/adapter/db"
	"github.com/aalperen0/syncognize/internal/application"
	"github.com/aalperen0/syncognize/internal/config"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("../.env"); err != nil {
		panic("failed load config" + err.Error())
	} else {
		fmt.Println(".env file loaded successfully from ../.env")
	}

	cfg, err := config.Load()
	if err != nil {
		panic("failed load config" + err.Error())
	}

	logger, err := cfg.NewLogger()
	if err != nil {
		panic("failed to create logger" + err.Error())
	}
	defer logger.Sync()

	dbConn, err := db.NewPostgres(ctx, cfg.Database, logger)
	if err != nil {
		logger.Fatal("database startup failed...", zap.Error(err))
	}
	defer dbConn.Close()

	// Test DB connection with simple query
	var result int
	if err := dbConn.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		logger.Fatal("database query test failed", zap.Error(err))
	}
	logger.Info("database query test successful", zap.Int("result", result))

	// Log connection pool stats
	stats := dbConn.Stats()
	logger.Info("database pool stats",
		zap.Int("open_connections", stats.OpenConnections),
		zap.Int("in_use", stats.InUse),
		zap.Int("idle", stats.Idle),
		zap.Int("max_open_conns", stats.MaxOpenConnections),
	)

	_ = &application.App{
		DB:     dbConn,
		Logger: logger,
	}

	logger.Info("syncognize starting...",
		zap.String("version", cfg.Service.Version),
		zap.String("environment", cfg.Service.Environment),
	)

}
