package config_test

import (
	"testing"

	"github.com/titoffon/merch-store/internal/config"
)

func TestLoadConfig(t *testing.T) {
    t.Setenv("PORT", "9090")
    t.Setenv("DB_HOST", "test_host")
	t.Setenv("LOG_LEVEL", "WARN")

    cfg := config.LoadConfig()

    if cfg.Port != "9090" {
        t.Errorf("expected Port=9090, got=%s", cfg.Port)
    }
    if cfg.DBHost != "test_host" {
        t.Errorf("expected DBHost=test_host, got=%s", cfg.DBHost)
    }

    if cfg.LogLevel != "WARN" {
        t.Errorf("expected default LOG_LEVEL=WARN, got=%s", cfg.LogLevel)
    }
}
