package config

import (
	"container/list"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(handler))
}

type Config struct {
	LogLevel string `json:"log_level"`
}

type ConfigWatcher struct {
	configPath   string
	callbacks    *list.List
	mu           sync.RWMutex
	cachedConfig *Config
}

func NewConfigWatcher(configPath string) *ConfigWatcher {
	slog.Debug("Creating new config watcher", "path", configPath)
	watcher := &ConfigWatcher{
		configPath: configPath,
		callbacks:  list.New(),
	}
	watcher.start()
	return watcher
}

func (cw *ConfigWatcher) AddCallback(callback func(*Config)) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.callbacks.PushBack(callback)
	slog.Debug("Added config callback", "total_callbacks", cw.callbacks.Len())
}

func (cw *ConfigWatcher) LoadConfig() *Config {
	cw.mu.RLock()
	if cw.cachedConfig != nil {
		config := cw.cachedConfig
		cw.mu.RUnlock()
		slog.Debug("Loaded config from cache", "log_level", config.LogLevel)
		return config
	}
	cw.mu.RUnlock()

	return cw.loadFromFile()
}

func (cw *ConfigWatcher) loadFromFile() *Config {
	slog.Debug("Loading config from file", "path", cw.configPath)

	data, err := os.ReadFile(cw.configPath)
	if err != nil {
		slog.Error("Failed to read config file", "error", err, "path", cw.configPath)
		config := &Config{LogLevel: "INFO"}
		cw.mu.Lock()
		cw.cachedConfig = config
		cw.mu.Unlock()
		slog.Debug("Using default config", "log_level", config.LogLevel)
		return config
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		slog.Error("Failed to parse config JSON", "error", err, "path", cw.configPath)
		defaultConfig := &Config{LogLevel: "INFO"}
		cw.mu.Lock()
		cw.cachedConfig = defaultConfig
		cw.mu.Unlock()
		slog.Debug("Using default config after JSON parse error", "log_level", defaultConfig.LogLevel)
		return defaultConfig
	}

	cw.mu.Lock()
	oldConfig := cw.cachedConfig
	cw.cachedConfig = &config
	cw.mu.Unlock()

	if oldConfig == nil || oldConfig.LogLevel != config.LogLevel {
		slog.Info("Config loaded successfully", "log_level", config.LogLevel, "path", cw.configPath)
	}

	slog.Debug("Successfully loaded config from file", "log_level", config.LogLevel, "path", cw.configPath)
	return &config
}

func (cw *ConfigWatcher) start() {
	slog.Debug("Starting config watcher", "path", cw.configPath)
	cw.loadFromFile()
	go cw.watchConfig()
}

func (cw *ConfigWatcher) watchConfig() {
	slog.Debug("Initializing file watcher", "path", cw.configPath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Failed to create file watcher", "error", err, "path", cw.configPath)
		return
	}
	defer watcher.Close()

	// Watch the directory containing the config file
	dir := filepath.Dir(cw.configPath)
	if dir == "." {
		dir = "./"
	}

	err = watcher.Add(dir)
	if err != nil {
		slog.Error("Failed to watch config directory", "error", err, "path", dir)
		return
	}

	slog.Debug("File watcher started successfully", "path", cw.configPath, "watching_dir", dir)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				slog.Debug("File watcher events channel closed", "path", cw.configPath)
				return
			}

			// Check if the event is for our config file
			if event.Name == cw.configPath {
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					slog.Debug("Config file change detected", "path", cw.configPath, "event", event.Op.String())
					slog.Info("Config file change detected", "path", cw.configPath)
					cw.reloadAndNotify()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				slog.Debug("File watcher errors channel closed", "path", cw.configPath)
				return
			}
			slog.Error("File watcher error", "error", err, "path", cw.configPath)
		}
	}
}

func (cw *ConfigWatcher) reloadAndNotify() {
	slog.Debug("Reloading config and notifying callbacks", "path", cw.configPath)

	config := cw.loadFromFile()
	cw.mu.RLock()
	callbackCount := cw.callbacks.Len()
	cw.mu.RUnlock()

	if callbackCount > 0 {
		slog.Info("Notifying config change callbacks", "callback_count", callbackCount, "log_level", config.LogLevel)
	}

	slog.Debug("Notifying callbacks of config change", "callback_count", callbackCount, "log_level", config.LogLevel)

	cw.mu.RLock()
	defer cw.mu.RUnlock()

	for element := cw.callbacks.Front(); element != nil; element = element.Next() {
		callback := element.Value.(func(*Config))
		callback(config)
	}

	slog.Debug("All callbacks notified", "callback_count", callbackCount)
}
