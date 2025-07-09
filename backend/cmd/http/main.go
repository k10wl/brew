package main

import (
	"log/slog"

	_ "brew/internal/util"
)

func main() {
	slog.Info("Starting brew HTTP server")
	slog.Info("HTTP server ready")
}
