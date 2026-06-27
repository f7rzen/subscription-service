package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/f7rzen/subscription-service/internal/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func NewPostgresDB(ctx context.Context, cfg config.Config, logger *slog.Logger) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	database, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	logger.Info("connected to postgres")

	return database, nil
}
