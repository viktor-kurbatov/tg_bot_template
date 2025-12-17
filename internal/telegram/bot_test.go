package telegram

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/viktor-kurbatov/tg_bot_template/pkg/logger"
)

func TestTelegram_logUpdateId_AddsUpdateIDToLogsViaContext(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&logger.Config{Writer: &buf, JSON: true})

	tg := &Telegram{logger: log.Logger}

	next := func(ctx context.Context, _ *bot.Bot, _ *models.Update) {
		tg.logger.InfoContext(ctx, "next called")
	}

	h := tg.logUpdateId(next)
	h(context.Background(), nil, &models.Update{ID: 42})

	out := buf.String()
	if !strings.Contains(out, "update_id") || !strings.Contains(out, "42") {
		t.Fatalf("expected log to contain update_id=42, got: %s", out)
	}
}

func TestTelegram_recoverBot_RecoversAndLogsPanic(t *testing.T) {
	var buf bytes.Buffer
	log := logger.New(&logger.Config{Writer: &buf, JSON: true})

	tg := &Telegram{logger: log.Logger}

	next := func(ctx context.Context, _ *bot.Bot, _ *models.Update) {
		panic("boom")
	}

	h := tg.recoverBot(next)
	// Should not panic.
	h(context.Background(), nil, &models.Update{ID: 1})

	out := buf.String()
	if !strings.Contains(out, "panic") || !strings.Contains(out, "boom") {
		t.Fatalf("expected log to contain panic + boom, got: %s", out)
	}
}
