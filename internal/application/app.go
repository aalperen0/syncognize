package application

import (
	"database/sql"

	"go.uber.org/zap"
)

type App struct {
	DB     *sql.DB
	Logger *zap.Logger
}
