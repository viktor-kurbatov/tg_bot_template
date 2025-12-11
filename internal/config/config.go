package config

import (
	"fmt"
	"os"
)

type Config struct {
	LogLevel    string // debug, info, warn, error
	BotToken    string
}

func Load() (*Config, error) {
	cfg := &Config{
		LogLevel:    os.Getenv("LOG_LEVEL"),
		BotToken:    os.Getenv("BOT_TOKEN"),
	}

	// Устанавливаем дефолтные значения, если переменные не заданы
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	// Проверяем обязательные переменные
	// Пока только BotToken важен для старта, DatabaseURL может быть позже
	if cfg.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}

	return cfg, nil
}