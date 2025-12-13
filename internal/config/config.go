package config

import (
	"fmt"
	"os"

	"github.com/viktor-kurbatov/tg_bot_template/pkg/logger"
)

type Config struct {
	Logger   *logger.Config
	Telegram *TelegramConfig
}

type TelegramConfig struct {
	Token string
}

func Load() (*Config, error) {
	cfg := &Config{
		Logger:   newLoggerConfig(),
		Telegram: &TelegramConfig{
			Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
		},
	}

	if cfg.Telegram.Token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	return cfg, nil
}

func newLoggerConfig() *logger.Config {
	cfg := &logger.Config{
		Level:     os.Getenv("LOG_LEVEL"),
		Output:    os.Getenv("LOG_OUTPUT"),
		File:      os.Getenv("LOG_FILE"),
		JSON:      os.Getenv("LOG_JSON") == "true",
		AddSource: os.Getenv("LOG_LEVEL") == "debug",
	}
	// Устанавливаем дефолтные значения, если переменные не заданы
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Output == "" {
		cfg.Output = "stdout"
	}
	return cfg
}
