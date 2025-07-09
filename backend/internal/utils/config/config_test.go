package config

import (
	"container/list"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func newTestConfigWatcher(configPath string) *ConfigWatcher {
	return &ConfigWatcher{
		configPath:     configPath,
		callbacks:      list.New(),
		watcherFactory: fsnotify.NewWatcher,
	}
}

func newTestConfigWatcherWithFactory(configPath string, factory WatcherFactory) *ConfigWatcher {
	return &ConfigWatcher{
		configPath:     configPath,
		callbacks:      list.New(),
		watcherFactory: factory,
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

func TestNewConfigWatcher(t *testing.T) {
	testFile := "test-new-watcher.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "DEBUG"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher := NewConfigWatcher(testFile)

	if watcher.configPath != testFile {
		t.Errorf("expected configPath %s, got %s", testFile, watcher.configPath)
	}

	if watcher.callbacks == nil {
		t.Error("expected callbacks list to be initialized")
	}

	config := watcher.LoadConfig()
	if config.LogLevel != "DEBUG" {
		t.Errorf("expected DEBUG from loaded config, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatcherDirectoryError(t *testing.T) {
	invalidPath := "/non/existent/path/config.json"

	watcher := newTestConfigWatcher(invalidPath)

	config := watcher.LoadConfig()
	if config.LogLevel != "INFO" {
		t.Errorf("expected default INFO for invalid path, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_EventChannelClosed(t *testing.T) {
	testFile := "test-event-closed.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher := newTestConfigWatcher(testFile)
	watcher.start()

	time.Sleep(50 * time.Millisecond)

	config := watcher.LoadConfig()
	if config.LogLevel != "INFO" {
		t.Errorf("expected INFO, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_ReloadAndNotify(t *testing.T) {
	testFile := "test-reload-notify.json"
	defer os.Remove(testFile)

	watcher := newTestConfigWatcher(testFile)

	var mu sync.Mutex
	var notifiedConfig *Config

	callback := func(cfg *Config) {
		mu.Lock()
		notifiedConfig = cfg
		mu.Unlock()
	}

	watcher.AddCallback(callback)

	testConfig := Config{LogLevel: "ERROR"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher.reloadAndNotify()

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if notifiedConfig == nil {
		t.Error("expected callback to be called")
		return
	}

	if notifiedConfig.LogLevel != "ERROR" {
		t.Errorf("expected ERROR from reloadAndNotify, got %s", notifiedConfig.LogLevel)
	}
}

func TestConfigWatcher_WatchConfigEventHandling(t *testing.T) {
	testFile := "test-watch-events.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher := newTestConfigWatcher(testFile)

	var mu sync.Mutex
	var callbackCalled bool

	callback := func(cfg *Config) {
		mu.Lock()
		callbackCalled = true
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

	if !callbackCalled {
		t.Error("expected callback to be called for file change")
	}

	config := watcher.LoadConfig()
	if config.LogLevel != "DEBUG" {
		t.Errorf("expected DEBUG after file change, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatchConfigInvalidDirectory(t *testing.T) {
	invalidPath := "/non/existent/dir/config.json"

	watcher := newTestConfigWatcher(invalidPath)
	watcher.start()

	time.Sleep(50 * time.Millisecond)

	config := watcher.LoadConfig()
	if config.LogLevel != "INFO" {
		t.Errorf("expected default INFO for invalid directory, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatchConfigCreateEvent(t *testing.T) {
	testFile := "test-create-event.json"
	defer os.Remove(testFile)

	watcher := newTestConfigWatcher(testFile)

	var mu sync.Mutex
	var callbackCalled bool

	callback := func(cfg *Config) {
		mu.Lock()
		callbackCalled = true
		mu.Unlock()
	}

	watcher.AddCallback(callback)
	watcher.start()

	time.Sleep(50 * time.Millisecond)

	testConfig := Config{LogLevel: "WARN"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if !callbackCalled {
		t.Error("expected callback to be called for file creation")
	}

	config := watcher.LoadConfig()
	if config.LogLevel != "WARN" {
		t.Errorf("expected WARN after file creation, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatchConfigDifferentFile(t *testing.T) {
	testFile := "test-different-file.json"
	otherFile := "test-other-file.json"
	defer os.Remove(testFile)
	defer os.Remove(otherFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher := newTestConfigWatcher(testFile)

	var mu sync.Mutex
	var callbackCalled bool

	callback := func(cfg *Config) {
		mu.Lock()
		callbackCalled = true
		mu.Unlock()
	}

	watcher.AddCallback(callback)
	watcher.start()

	time.Sleep(50 * time.Millisecond)

	otherConfig := Config{LogLevel: "ERROR"}
	data, _ = json.Marshal(otherConfig)
	os.WriteFile(otherFile, data, 0644)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if callbackCalled {
		t.Error("expected callback NOT to be called for different file")
	}

	config := watcher.LoadConfig()
	if config.LogLevel != "INFO" {
		t.Errorf("expected INFO unchanged, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatchConfigGracefulDegradation(t *testing.T) {
	testFile := "test-graceful-degradation.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "WARN"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	watcher := newTestConfigWatcher(testFile)

	go watcher.watchConfig()

	time.Sleep(50 * time.Millisecond)

	config := watcher.LoadConfig()
	if config.LogLevel != "WARN" {
		t.Errorf("expected WARN, got %s", config.LogLevel)
	}

	os.Remove(testFile)

	updatedConfig := Config{LogLevel: "ERROR"}
	data, _ = json.Marshal(updatedConfig)
	os.WriteFile(testFile, data, 0644)

	time.Sleep(100 * time.Millisecond)

	config = watcher.LoadConfig()
	if config.LogLevel != "ERROR" {
		t.Errorf("expected ERROR after recreation, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatcherCreationFailure(t *testing.T) {
	testFile := "test-watcher-creation-failure.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "DEBUG"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	failingFactory := func() (*fsnotify.Watcher, error) {
		return nil, errors.New("simulated watcher creation failure")
	}

	watcher := newTestConfigWatcherWithFactory(testFile, failingFactory)

	watcher.loadFromFile()

	go watcher.watchConfig()

	time.Sleep(100 * time.Millisecond)

	config := watcher.LoadConfig()
	if config.LogLevel != "DEBUG" {
		t.Errorf("expected DEBUG despite watcher failure, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatchConfigErrorsChannelClosed(t *testing.T) {
	testFile := "test-errors-channel-closed.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	mockWatcher := newMockWatcher()

	factory := func() (*fsnotify.Watcher, error) {
		return &mockWatcher.Watcher, nil
	}

	watcher := newTestConfigWatcherWithFactory(testFile, factory)

	done := make(chan bool)
	go func() {
		watcher.watchConfig()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)

	close(mockWatcher.errors)

	select {
	case <-done:
		// Expected - watchConfig should return when errors channel is closed
	case <-time.After(time.Second):
		t.Error("expected watchConfig to return when errors channel closed")
	}
}

func TestConfigWatcher_WatchConfigErrorReceived(t *testing.T) {
	testFile := "test-error-received.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	mockWatcher := newMockWatcher()

	factory := func() (*fsnotify.Watcher, error) {
		return &mockWatcher.Watcher, nil
	}

	watcher := newTestConfigWatcherWithFactory(testFile, factory)

	go watcher.watchConfig()

	time.Sleep(50 * time.Millisecond)

	testError := errors.New("test watcher error")
	mockWatcher.errors <- testError

	time.Sleep(50 * time.Millisecond)

	config := watcher.LoadConfig()
	if config.LogLevel != "INFO" {
		t.Errorf("expected INFO after error, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatchConfigEventsChannelClosed(t *testing.T) {
	testFile := "test-events-channel-closed.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	mockWatcher := newMockWatcher()

	factory := func() (*fsnotify.Watcher, error) {
		return &mockWatcher.Watcher, nil
	}

	watcher := newTestConfigWatcherWithFactory(testFile, factory)

	done := make(chan bool)
	go func() {
		watcher.watchConfig()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)

	close(mockWatcher.events)

	select {
	case <-done:
		// Expected - watchConfig should return when events channel is closed
	case <-time.After(time.Second):
		t.Error("expected watchConfig to return when events channel closed")
	}
}

func TestConfigWatcher_WatchConfigWriteEvent(t *testing.T) {
	testFile := "test-write-event.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	mockWatcher := newMockWatcher()

	factory := func() (*fsnotify.Watcher, error) {
		return &mockWatcher.Watcher, nil
	}

	watcher := newTestConfigWatcherWithFactory(testFile, factory)

	var mu sync.Mutex
	var callbackTriggered bool

	callback := func(cfg *Config) {
		mu.Lock()
		callbackTriggered = true
		mu.Unlock()
	}

	watcher.AddCallback(callback)

	go watcher.watchConfig()

	time.Sleep(50 * time.Millisecond)

	updatedConfig := Config{LogLevel: "DEBUG"}
	data, _ = json.Marshal(updatedConfig)
	os.WriteFile(testFile, data, 0644)

	writeEvent := fsnotify.Event{
		Name: testFile,
		Op:   fsnotify.Write,
	}

	mockWatcher.events <- writeEvent

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if !callbackTriggered {
		t.Error("expected callback to be triggered for Write event")
	}

	config := watcher.LoadConfig()
	if config.LogLevel != "DEBUG" {
		t.Errorf("expected DEBUG after Write event, got %s", config.LogLevel)
	}
}

func TestConfigWatcher_WatchConfigIgnoreOtherOperations(t *testing.T) {
	testFile := "test-ignore-other-ops.json"
	defer os.Remove(testFile)

	testConfig := Config{LogLevel: "INFO"}
	data, _ := json.Marshal(testConfig)
	os.WriteFile(testFile, data, 0644)

	mockWatcher := newMockWatcher()

	factory := func() (*fsnotify.Watcher, error) {
		return &mockWatcher.Watcher, nil
	}

	watcher := newTestConfigWatcherWithFactory(testFile, factory)

	var mu sync.Mutex
	var callbackTriggered bool

	callback := func(cfg *Config) {
		mu.Lock()
		callbackTriggered = true
		mu.Unlock()
	}

	watcher.AddCallback(callback)

	go watcher.watchConfig()

	time.Sleep(50 * time.Millisecond)

	removeEvent := fsnotify.Event{
		Name: testFile,
		Op:   fsnotify.Remove,
	}

	mockWatcher.events <- removeEvent

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if callbackTriggered {
		t.Error("expected callback NOT to be triggered for Remove event")
	}

	config := watcher.LoadConfig()
	if config.LogLevel != "INFO" {
		t.Errorf("expected INFO unchanged after Remove event, got %s", config.LogLevel)
	}
}

type MockWatcher struct {
	fsnotify.Watcher
	events chan fsnotify.Event
	errors chan error
}

func newMockWatcher() *MockWatcher {
	mock := &MockWatcher{
		events: make(chan fsnotify.Event, 1),
		errors: make(chan error, 1),
	}
	// Create a real watcher to embed
	realWatcher, _ := fsnotify.NewWatcher()
	mock.Watcher = *realWatcher
	// Replace the channels with our mock ones
	mock.Watcher.Events = mock.events
	mock.Watcher.Errors = mock.errors
	return mock
}

func (m *MockWatcher) Add(name string) error {
	return nil
}

func (m *MockWatcher) Remove(name string) error {
	return nil
}

func (m *MockWatcher) Close() error {
	return nil
}
