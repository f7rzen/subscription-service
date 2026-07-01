package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/f7rzen/subscription-service/internal/config"
	"github.com/f7rzen/subscription-service/internal/db"
	"github.com/f7rzen/subscription-service/internal/handler"
	"github.com/f7rzen/subscription-service/internal/repository"
	"github.com/f7rzen/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx := context.Background()

	database, err := db.NewPostgresDB(ctx, cfg, logger)
	if err != nil {
		logger.Error(
			"failed to connect to postgres",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer database.Close()

	subscriptionRepository := repository.NewSubscriptionRepository(database)
	subscriptionService := service.NewSubscriptionService(subscriptionRepository, logger)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, logger)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.Create)
		}
	}

	logger.Info(
		"starting server",
		slog.String("port", cfg.AppPort),
	)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		logger.Error(
			"failed to start server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
