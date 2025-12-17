package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("Default values", func(t *testing.T) {
		t.Setenv("TELEGRAM_BOT_TOKEN", "test_token")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Logger.Level != "info" {
			t.Errorf("Expected default LogLevel 'info', got '%s'", cfg.Logger.Level)
		}
		if cfg.Logger.Output != "stdout" {
			t.Errorf("Expected default LogOutput 'stdout', got '%s'", cfg.Logger.Output)
		}
	})

	t.Run("Custom values", func(t *testing.T) {
		t.Setenv("TELEGRAM_BOT_TOKEN", "test_token")
		t.Setenv("LOG_LEVEL", "debug")
		t.Setenv("LOG_OUTPUT", "file")
		t.Setenv("LOG_FILE", "app.log")
		t.Setenv("LOG_JSON", "true")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if cfg.Logger.Level != "debug" {
			t.Errorf("Expected LogLevel 'debug', got '%s'", cfg.Logger.Level)
		}
		if cfg.Logger.Output != "file" {
			t.Errorf("Expected LogOutput 'file', got '%s'", cfg.Logger.Output)
		}
		if cfg.Logger.File != "app.log" {
			t.Errorf("Expected LogFile 'app.log', got '%s'", cfg.Logger.File)
		}
		if !cfg.Logger.JSON {
			t.Errorf("Expected JSON true, got false")
		}
		// AddSource should be true for debug level
		if !cfg.Logger.AddSource {
			t.Errorf("Expected AddSource true for debug level, got false")
		}
	})

	t.Run("Missing required vars", func(t *testing.T) {
		t.Setenv("TELEGRAM_BOT_TOKEN", "")

		_, err := Load()
		if err == nil {
			t.Error("Expected error for missing TELEGRAM_BOT_TOKEN, got nil")
		}
	})
}
