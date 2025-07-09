package main

import (
	"brew/internal/utils/logger"
	"time"
)

func main() {
	logger.Info("Starting brew HTTP server")
	logger.Info("HTTP server ready")

	counter := 0
	for {
		counter++

		logger.Info("Application running", "iteration", counter, "status", "healthy")
		logger.Debug("Debug message for troubleshooting", "iteration", counter, "timestamp", time.Now().Format(time.RFC3339))
		logger.Error("Simulated error for testing", "iteration", counter, "error_type", "test_error")

		time.Sleep(1 * time.Second)
	}
}
