package main

import (
	"log/slog"
	"os"
	"ticket-office/internal/config"
	"ticket-office/internal/storage/sqlite"
)

const (
	devLogLevel   = "dev"
	localLogLevel = "local"
)

func main() {
	cfg := config.MustLoad()

	log := NewLog(cfg.LogLevel)
	log.Info("logger was initialized")

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error(err.Error())
	}

	_ = storage
}

func NewLog(logLevel string) *slog.Logger {
	var log *slog.Logger

	switch logLevel {
	case devLogLevel:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case localLogLevel:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
