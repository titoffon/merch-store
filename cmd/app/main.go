package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/titoffon/merch-store/internal/config"
	"github.com/titoffon/merch-store/internal/db"
	"github.com/titoffon/merch-store/internal/delivery/routes"
	"github.com/titoffon/merch-store/pkg/logger"
)

func main(){

	cfg := config.LoadConfig()
	

	logger.InitGlobalLogger(cfg.LogLevel)

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	ctx := context.Background()

	pool, err := db.New(ctx, connectionString)
	if err != nil {
		slog.Error("failed to coыnnect to database: %v", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.DBPool.Close()

	r := routes.NewRouter(pool)

	slog.Info("Starting server", slog.String("address", cfg.Port))
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}

}