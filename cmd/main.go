package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/server/handlers/url/redirect"
	"url-shortener/internal/server/handlers/url/save"
	middleware_logger "url-shortener/internal/server/middleware/logger"
	"url-shortener/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("startup url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug log are enabled")

	storage, err := sqlite.NewStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("err", err))
		os.Exit(1)
	}
	println(storage)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware_logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HttpServer.User: cfg.HttpServer.Password,
		}))
		r.Post("/", save.New(log, storage))
		// TODO DELETE /url/{id}
	})
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.HttpServer.Address))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true}))
	}

	return logger
}
