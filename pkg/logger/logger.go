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
	Level     string    // debug, info, warn, error
	Output    string    // stdout, stderr, file
	File      string    // path for file output
	JSON      bool      // use JSON format
	AddSource bool      // add source file info
	Writer    io.Writer // custom writer (overrides Output, useful for testing)
}

// Logger wraps slog.Logger with Close() support for file writers
type Logger struct {
	*slog.Logger
	closer io.Closer
}

// Close closes the underlying writer if it implements io.Closer
func (l *Logger) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

type ctxKey string

const slogAttrs ctxKey = "slog_attrs"

// New creates a new structured logger
func New(cfg *Config) *Logger {
	opts := &slog.HandlerOptions{
		Level:     parseLogLevel(cfg.Level),
		AddSource: cfg.AddSource,
	}

	output, closer := parseLogOutput(cfg)

	var handler slog.Handler
	if cfg.JSON {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	return &Logger{
		Logger: slog.New(&ctxHandler{Handler: handler}),
		closer: closer,
	}
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
func parseLogOutput(cfg *Config) (io.Writer, io.Closer) {
	// Custom writer takes priority (useful for testing)
	if cfg.Writer != nil {
		return cfg.Writer, nil
	}

	switch cfg.Output {
	case "stdout", "":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	case "file":
		if cfg.File == "" {
			cfg.File = "logs/app.log"
		}
		lj := &lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   false,
		}
		return lj, lj
	default:
		return os.Stdout, nil
	}
}
