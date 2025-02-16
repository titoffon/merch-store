package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/titoffon/merch-store/internal/config"
	"github.com/titoffon/merch-store/internal/db"
	"github.com/titoffon/merch-store/internal/delivery/routes"
	"github.com/titoffon/merch-store/pkg/logger"
)

func Run(cfg *config.Config ) error{
	logger.InitGlobalLogger(cfg.LogLevel)

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	ctx := context.Background()

	dal, err := db.New(ctx, connectionString)
	if err != nil {
		slog.Error("failed to coыnnect to database: %v", slog.String("error", err.Error()))
		return err
	}
	defer dal.DBPool.Close()

	r := routes.NewRouter(dal)

	slog.Info("Starting server", slog.String("address", cfg.Port))
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
		return err
	}
	return nil
}