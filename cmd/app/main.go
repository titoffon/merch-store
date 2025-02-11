package main

import (
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

	pool, err := db.InitDB(cfg)
	if err != nil {
		slog.Error("failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer pool.Close()

	r := routes.NewRouter(pool)

	slog.Info("Starting server", slog.String("address", cfg.Port))
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}

}