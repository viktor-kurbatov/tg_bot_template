package telegram

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/viktor-kurbatov/tg_bot_template/internal/config"
	"github.com/viktor-kurbatov/tg_bot_template/pkg/logger"
)

type Telegram struct {
	telegramBot *bot.Bot
	logger      *slog.Logger
}

func NewBot(cfg *config.Config, logger *slog.Logger) (*Telegram, error) {
	logger = logger.With("component", "telegram")
	telegram := &Telegram{logger: logger}
	telegramBot, err := bot.New(cfg.Telegram.Token,
		bot.WithDefaultHandler(telegram.HandleUpdate),
		bot.WithMiddlewares(telegram.recoverBot, telegram.logUpdateId),
	)
	if err != nil {
		return nil, err
	}

	telegram.telegramBot = telegramBot

	return telegram, nil
}

func (b *Telegram) Start(ctx context.Context) {
	b.logger.Info("Starting telegram bot")
	b.telegramBot.Start(ctx)
}

func (b *Telegram) HandleUpdate(ctx context.Context, bot *bot.Bot, update *models.Update) {
	b.logger.InfoContext(ctx, "Handling update", "update", update)
}

func (b *Telegram) recoverBot(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		defer func() {
			if r := recover(); r != nil {
				b.logger.ErrorContext(ctx, "panic", "error", r, "stack", string(debug.Stack()))
			}
		}()
		next(ctx, bot, update)
	}
}

func (b *Telegram) logUpdateId(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		ctx = logger.WithAttr(ctx, "update_id", update.ID)
		next(ctx, bot, update)
	}
}
