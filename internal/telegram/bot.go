package telegram

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/viktor-kurbatov/tg_bot_template/internal/config"
)

type Telegram struct {
	telegramBot *bot.Bot
	logger      *slog.Logger
}

func NewBot(cfg *config.Config, logger *slog.Logger) (*Telegram, error) {
	logger = logger.With("component", "telegram")
	telegramBot, err := bot.New(cfg.Telegram.Token, bot.WithDebug())
	if err != nil {
		return nil, err
	}
	return &Telegram{telegramBot: telegramBot, logger: logger}, nil
}

func (b *Telegram) Start(ctx context.Context) error {
	b.logger.Info("Starting telegram bot")
	b.telegramBot.Start(ctx)
	return nil
}
