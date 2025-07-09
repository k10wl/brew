package config

import (
	"container/list"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"
)

func newTestConfigWatcher(configPath string) *ConfigWatcher {
	return &ConfigWatcher{
		configPath: configPath,
		callbacks:  list.New(),
	}
}

func TestConfigWatcher_LoadConfig(t *testing.T) {
	testFile := "test-config.json"
	defer os.Remove(testFile)

	watcher := newTestConfigWatcher(testFile)

	t.Run("loads default config when file doesn't exist", func(t *testing.T) {
		config := watcher.LoadConfig()
		if config.LogLevel != "INFO" {
			t.Errorf("expected default log level INFO, got %s", config.LogLevel)
		}
	})

	t.Run("loads config from file", func(t *testing.T) {
		testConfig := Config{LogLevel: "DEBUG"}
		data, _ := json.Marshal(testConfig)
		os.WriteFile(testFile, data, 0644)

		watcher.cachedConfig = nil
		config := watcher.LoadConfig()
		if config.LogLevel != "DEBUG" {
			t.Errorf("expected DEBUG, got %s", config.LogLevel)
		}
	})

	t.Run("handles invalid JSON", func(t *testing.T) {
		os.WriteFile(testFile, []byte("invalid json"), 0644)

		watcher.cachedConfig = nil
		config := watcher.LoadConfig()
		if config.LogLevel != "INFO" {
			t.Errorf("expected default INFO for invalid JSON, got %s", config.LogLevel)
		}
	})
}

func TestConfigWatcher_SingletonCache(t *testing.T) {
	testFile := "test-singleton.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "WARN"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher := newTestConfigWatcher(testFile)

	config1 := watcher.LoadConfig()
	config2 := watcher.LoadConfig()

	if config1 != config2 {
		t.Error("expected same config instance from cache")
	}

	if config1.LogLevel != "WARN" {
		t.Errorf("expected WARN, got %s", config1.LogLevel)
	}
}

func TestConfigWatcher_MultipleCallbacks(t *testing.T) {
	testFile := "test-callbacks.json"
	defer os.Remove(testFile)

	watcher := newTestConfigWatcher(testFile)

	var mu sync.Mutex
	var calls []string

	callback1 := func(cfg *Config) {
		mu.Lock()
		calls = append(calls, "callback1:"+cfg.LogLevel)
		mu.Unlock()
	}

	callback2 := func(cfg *Config) {
		mu.Lock()
		calls = append(calls, "callback2:"+cfg.LogLevel)
		mu.Unlock()
	}

	watcher.AddCallback(callback1)
	watcher.AddCallback(callback2)

	testConfig := Config{LogLevel: "ERROR"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher.reloadAndNotify()

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(calls) != 2 {
		t.Errorf("expected 2 callback calls, got %d", len(calls))
	}

	expectedCalls := []string{"callback1:ERROR", "callback2:ERROR"}
	for i, expected := range expectedCalls {
		if i >= len(calls) || calls[i] != expected {
			t.Errorf("expected call %d to be %s, got %s", i, expected, calls[i])
		}
	}
}

func TestConfigWatcher_FileWatch(t *testing.T) {
	testFile := "test-watch.json"
	defer os.Remove(testFile)

	initialConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(initialConfig)
	os.WriteFile(testFile, data, 0644)

	watcher := newTestConfigWatcher(testFile)

	var mu sync.Mutex
	var receivedConfig *Config

	callback := func(cfg *Config) {
		mu.Lock()
		receivedConfig = cfg
		mu.Unlock()
	}

	watcher.AddCallback(callback)
	watcher.start()

	time.Sleep(50 * time.Millisecond)

	updatedConfig := Config{LogLevel: "DEBUG"}
	data, _ = json.Marshal(updatedConfig)
	os.WriteFile(testFile, data, 0644)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if receivedConfig == nil {
		t.Error("expected callback to be called")
		return
	}

	if receivedConfig.LogLevel != "DEBUG" {
		t.Errorf("expected DEBUG from file watch, got %s", receivedConfig.LogLevel)
	}
}

func TestConfigWatcher_ConcurrentAccess(t *testing.T) {
	testFile := "test-concurrent.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher := newTestConfigWatcher(testFile)

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			config := watcher.LoadConfig()
			if config.LogLevel != "INFO" {
				t.Errorf("expected INFO, got %s", config.LogLevel)
			}
		}()
	}

	wg.Wait()
}

func TestConfigWatcher_CallbackQueue(t *testing.T) {
	testFile := "test-queue.json"
	defer os.Remove(testFile)

	watcher := newTestConfigWatcher(testFile)

	var mu sync.Mutex
	var callOrder []int

	for i := 0; i < 5; i++ {
		id := i
		callback := func(cfg *Config) {
			mu.Lock()
			callOrder = append(callOrder, id)
			mu.Unlock()
		}
		watcher.AddCallback(callback)
	}

	testConfig := Config{LogLevel: "WARN"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher.reloadAndNotify()

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(callOrder) != 5 {
		t.Errorf("expected 5 callbacks, got %d", len(callOrder))
	}

	for i, id := range callOrder {
		if id != i {
			t.Errorf("expected callback %d, got %d", i, id)
		}
	}
}
