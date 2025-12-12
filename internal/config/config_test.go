package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Сохраняем текущие env vars, чтобы восстановить после теста
	originalToken := os.Getenv("BOT_TOKEN")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	originalLogOutput := os.Getenv("LOG_OUTPUT")

	defer func() {
		os.Setenv("BOT_TOKEN", originalToken)
		os.Setenv("LOG_LEVEL", originalLogLevel)
		os.Setenv("LOG_OUTPUT", originalLogOutput)
	}()

	t.Run("Default values", func(t *testing.T) {
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_OUTPUT")
		os.Setenv("BOT_TOKEN", "test_token")

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
		os.Setenv("BOT_TOKEN", "test_token")
		os.Setenv("LOG_LEVEL", "debug")
		os.Setenv("LOG_OUTPUT", "file")
		os.Setenv("LOG_FILE", "app.log")
		os.Setenv("LOG_JSON", "true")

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
		os.Unsetenv("BOT_TOKEN")

		_, err := Load()
		if err == nil {
			t.Error("Expected error for missing BOT_TOKEN, got nil")
		}
	})
}
