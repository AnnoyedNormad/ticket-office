package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"os"
	"ticket-office/internal/config"
	"ticket-office/internal/http-server/handlers/event/saveevent"
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

	e := echo.New()

	e.Use(middleware.Recover(), middleware.Logger())

	eg := e.Group("/event")

	eg.Use(middleware.BasicAuth(func(user, pass string, e echo.Context) (bool, error) {
		if user == "home" && pass == "home" {
			return true, nil
		}
		return false, nil
	}))
	eg.POST("/save", saveevent.New(storage))

	//e.POST("/", saveuserstickets.New(storage))

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      e,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	server.ListenAndServe()
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
