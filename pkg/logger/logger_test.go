package logger

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"debug", "debug", slog.LevelDebug},
		{"info", "info", slog.LevelInfo},
		{"warn", "warn", slog.LevelWarn},
		{"error", "error", slog.LevelError},
		{"uppercase", "DEBUG", slog.LevelDebug},
		{"mixed case", "WaRn", slog.LevelWarn},
		{"with spaces", "  info  ", slog.LevelInfo},
		{"with comment", "debug#info", slog.LevelDebug},
		{"with spaces and comment", " debug # info ", slog.LevelDebug},
		{"unknown", "unknown", slog.LevelInfo},
		{"empty", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLogLevel(tt.input)
			if got != tt.expected {
				t.Errorf("parseLogLevel(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseLogOutput(t *testing.T) {
	var buf bytes.Buffer

	tests := []struct {
		name       string
		cfg        *Config
		wantWriter any
		wantCloser bool
	}{
		{"custom writer overrides output", &Config{Writer: &buf, Output: "file"}, &buf, false},
		{"stdout", &Config{Output: "stdout"}, os.Stdout, false},
		{"stderr", &Config{Output: "stderr"}, os.Stderr, false},
		{"empty defaults to stdout", &Config{Output: ""}, os.Stdout, false},
		{"unknown defaults to stdout", &Config{Output: "unknown"}, os.Stdout, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, closer := parseLogOutput(tt.cfg)
			if w != tt.wantWriter {
				t.Errorf("writer = %v, want %v", w, tt.wantWriter)
			}
			if (closer != nil) != tt.wantCloser {
				t.Errorf("closer = %v, wantCloser = %v", closer, tt.wantCloser)
			}
		})
	}
}

func TestParseLogOutput_File(t *testing.T) {
	tmpfile, err := os.CreateTemp(".", "test_log_*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmpPath := tmpfile.Name()
	tmpfile.Close()

	cfg := &Config{Output: "file", File: tmpPath}
	log := New(cfg)

	t.Cleanup(func() {
		log.Close()
		if err := os.Remove(tmpPath); err != nil {
			t.Fatalf("failed to remove log file: %v", err)
		}
	})

	testMsg := "hello from test logger"
	log.Info(testMsg)

	content, err := os.ReadFile(tmpPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), testMsg) {
		t.Errorf("expected log to contain %q, got %q", testMsg, content)
	}
}

func TestWithAttrs(t *testing.T) {
	ctx := context.Background()
	ctx = WithAttr(ctx, "key1", "value1")
	ctx = WithAttrs(ctx, slog.String("key2", "value2"))

	attrs, ok := ctx.Value(slogAttrs).([]slog.Attr)
	if !ok {
		t.Fatal("Context does not contain slogAttrs")
	}

	if len(attrs) != 2 {
		t.Errorf("Expected 2 attributes, got %d", len(attrs))
	}

	found := make(map[string]any)
	for _, a := range attrs {
		found[a.Key] = a.Value.Any()
	}

	if found["key1"] != "value1" {
		t.Errorf("Expected key1=value1, got %v", found["key1"])
	}
	if found["key2"] != "value2" {
		t.Errorf("Expected key2=value2, got %v", found["key2"])
	}
}
