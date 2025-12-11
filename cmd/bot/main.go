package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

)

func main() {
	fmt.Println("Starting bot app...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Create application
	application, err := NewApp(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create application: %v\n", err)
		os.Exit(1)
	}

	// Start application
	application.logger.Debug("application starting")
	if err := application.Start(ctx); err != nil {
		application.logger.Error("application error", "error", err)
		os.Exit(1)
	}

	application.logger.Info("application stopped")
}

type app struct {
	logger *slog.Logger
}

func NewApp(ctx context.Context) (*app, error) {
	logger, err := NewLogger()
	if err != nil {
		return nil, err
	}
	return &app{logger: logger}, nil
}

func (a *app) Start(ctx context.Context) error {
	a.logger.Debug("application starting")
	return nil
}

func NewLogger() (*slog.Logger, error) {
	var level slog.Level
	switch os.Getenv("LOG_LEVEL") {
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