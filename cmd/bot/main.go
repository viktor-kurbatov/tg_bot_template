package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/viktor-kurbatov/tg_bot_template/internal/config"
	"github.com/viktor-kurbatov/tg_bot_template/internal/telegram"
	"github.com/viktor-kurbatov/tg_bot_template/pkg/logger"
)

var (
	newLogger      = logger.New
	newTelegramBot = telegram.NewBot
)

func main() {
	fmt.Println("go running bot app...")

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
	defer application.Close()

	if err := application.Start(ctx); err != nil {
		application.logger.Error("application error", "error", err)
		os.Exit(1)
	}

	application.logger.Info("application stopped")
}

type App struct {
	logger      *logger.Logger
	telegramBot *telegram.Telegram
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	log := newLogger(cfg.Logger)
	bot, err := newTelegramBot(cfg, log.Logger)
	if err != nil {
		log.Close() // cleanup on failure
		return nil, err
	}
	return &App{logger: log, telegramBot: bot}, nil
}

func (a *App) Close() error {
	return a.logger.Close()
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Debug("application starting")
	a.telegramBot.Start(ctx)
	return nil
}
