package util

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Config struct {
	LogLevel string `json:"log_level"`
}

var (
	currentHandler *slog.JSONHandler
	handlerMu      sync.RWMutex
	configPath     = "http-config.json"
)

func init() {
	initLogger()
	go watchConfig()
}

func initLogger() {
	config := loadConfig()
	level := parseLogLevel(config.LogLevel)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	handlerMu.Lock()
	currentHandler = handler
	handlerMu.Unlock()

	slog.SetDefault(slog.New(handler))
}

func loadConfig() *Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{LogLevel: "INFO"}
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return &Config{LogLevel: "INFO"}
	}

	return &config
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

func watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Failed to create file watcher", "error", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(configPath)
	if err != nil {
		slog.Error("Failed to watch config file", "error", err, "path", configPath)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				updateLogger()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error("File watcher error", "error", err)
		}
	}
}

func updateLogger() {
	config := loadConfig()
	level := parseLogLevel(config.LogLevel)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	handlerMu.Lock()
	currentHandler = handler
	handlerMu.Unlock()

	slog.SetDefault(slog.New(handler))
	slog.Info("Log level updated", "new_level", config.LogLevel)
}

// Debug logs a debug message with key-value pairs
func Debug(msg string, keysAndValues ...any) {
	slog.Debug(msg, keysAndValues...)
}

// Info logs an info message with key-value pairs
func Info(msg string, keysAndValues ...any) {
	slog.Info(msg, keysAndValues...)
}

// Warn logs a warning message with key-value pairs
func Warn(msg string, keysAndValues ...any) {
	slog.Warn(msg, keysAndValues...)
}

// Error logs an error message with key-value pairs
func Error(msg string, keysAndValues ...any) {
	slog.Error(msg, keysAndValues...)
}
