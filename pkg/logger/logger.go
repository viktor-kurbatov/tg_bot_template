// pkg/logger/logger.go
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Config represents logger configuration
type Config struct {
	Level     string // debug, info, warn, error
	Output    string // stdout, file
	File      string
	JSON      bool
	AddSource bool
}

type ctxKey string

const slogAttrs ctxKey = "slog_attrs"

// New creates a new structured logger
func New(cfg *Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     parseLogLevel(cfg.Level),
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	if cfg.JSON {
		handler = slog.NewJSONHandler(parseLogOutput(cfg), opts)
	} else {
		handler = slog.NewTextHandler(parseLogOutput(cfg), opts)
	}

	return slog.New(&ctxHandler{Handler: handler})
}

// ctxHandler enriches log records with context attributes
type ctxHandler struct {
	slog.Handler
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogAttrs).([]slog.Attr); ok {
		r.AddAttrs(attrs...)
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns a new context with additional log attributes
func WithAttrs(parent context.Context, attrs ...slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	existing, _ := parent.Value(slogAttrs).([]slog.Attr)
	combined := append(existing, attrs...)

	return context.WithValue(parent, slogAttrs, combined)
}

// WithAttr returns a new context with a single attribute
func WithAttr(parent context.Context, key string, value any) context.Context {
	return WithAttrs(parent, slog.Any(key, value))
}

// parseLogLevel parses log level from string
func parseLogLevel(level string) slog.Level {
	level = strings.ToLower(level)
	level = strings.Split(level, "#")[0]
	level = strings.TrimSpace(level)
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// parseLogOutput determines output writer from environment
func parseLogOutput(cfg *Config) io.Writer {
	switch cfg.Output {
	case "stdout", "":
		return os.Stdout
	case "stderr":
		return os.Stderr
	case "file":
		if cfg.File == "" {
			cfg.File = "logs/app.log"
		}
		return &lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   false,
		}
	default:
		return os.Stdout
	}
}
