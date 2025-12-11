package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/viktor-kurbatov/tg_bot_template/internal/config"
)

func main() {
	fmt.Println("Starting bot app...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Create application
	application, err := NewApp(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create application: %v\n", err)
		os.Exit(1)
	}

	// Start application
	application.logger.Debug("The app is going to start.")
	if err := application.Start(ctx); err != nil {
		application.logger.Error("application error", "error", err)
		os.Exit(1)
	}

	application.logger.Info("application stopped")
}

type app struct {
	logger *slog.Logger
}

func NewApp(ctx context.Context, cfg *config.Config) (*app, error) {
	logger, err := NewLogger(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	return &app{logger: logger}, nil
}

func (a *app) Start(ctx context.Context) error {
	a.logger.Debug("application starting")
	return nil
}

func NewLogger(levelStr string) (*slog.Logger, error) {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	logger.Info("logger initialized", "level", level)
	return logger, nil
}
