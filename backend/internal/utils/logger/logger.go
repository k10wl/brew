package logger

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"brew/internal/utils/config"
)

var (
	currentHandler *slog.JSONHandler
	handlerMu      sync.RWMutex
	configWatcher  *config.ConfigWatcher
)

func init() {
	configWatcher = config.NewConfigWatcher("http-config.json")
	configWatcher.AddCallback(updateLogger)

	initialConfig := configWatcher.LoadConfig()
	updateLogger(initialConfig)
}

func updateLogger(cfg *config.Config) {
	level := parseLogLevel(cfg.LogLevel)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	handlerMu.Lock()
	currentHandler = handler
	handlerMu.Unlock()

	slog.SetDefault(slog.New(handler))
	slog.Info("Log level updated", "new_level", cfg.LogLevel)
}

func parseLogLevel(levelStr string) slog.Level {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Debug(msg string, keysAndValues ...any) {
	slog.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...any) {
	slog.Info(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...any) {
	slog.Warn(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...any) {
	slog.Error(msg, keysAndValues...)
}
