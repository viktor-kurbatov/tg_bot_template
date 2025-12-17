package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/viktor-kurbatov/tg_bot_template/internal/config"
	"github.com/viktor-kurbatov/tg_bot_template/internal/telegram"
	"github.com/viktor-kurbatov/tg_bot_template/pkg/logger"
)

func TestNewApp(t *testing.T) {
	origNewLogger := newLogger
	origNewTelegramBot := newTelegramBot
	t.Cleanup(func() {
		newLogger = origNewLogger
		newTelegramBot = origNewTelegramBot
	})

	cfg := &config.Config{
		Logger:   &logger.Config{Level: "info"},
		Telegram: &config.TelegramConfig{Token: "token"},
	}

	tests := []struct {
		name    string
		log     *logger.Logger
		bot     *telegram.Telegram
		botErr  error
		wantErr bool
	}{
		{
			name:    "success",
			log:     &logger.Logger{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))},
			bot:     &telegram.Telegram{},
			botErr:  nil,
			wantErr: false,
		},
		{
			name:    "bot error propagates",
			log:     &logger.Logger{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))},
			bot:     nil,
			botErr:  errors.New("bot creation failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newLogger = func(cfg *logger.Config) *logger.Logger {
				return tt.log
			}
			newTelegramBot = func(cfg *config.Config, log *slog.Logger) (*telegram.Telegram, error) {
				return tt.bot, tt.botErr
			}

			app, err := NewApp(context.Background(), cfg)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if app.logger != tt.log {
					t.Fatal("app.logger != expected logger")
				}
				if app.telegramBot != tt.bot {
					t.Fatal("app.telegramBot != expected bot")
				}
			}
		})
	}
}
