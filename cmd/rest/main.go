package main

import (
	"effective-mobile-test/internal/config"
	"effective-mobile-test/internal/db/postgresql"
	"effective-mobile-test/internal/http/handlers/v1"
	"effective-mobile-test/internal/usecases"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	db, err := postgres.New(log, cfg.DbPath)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	slp := postgres.NewSongLibrary(db)
	sluc := usecases.NewSongLibrary(slp, log)
	handlers.NewRouter(log, r, sluc)

	server := &http.Server{
		Addr:         cfg.HttpAddr,
		Handler:      r,
		ReadTimeout:  cfg.HttpReadTimeout,
		WriteTimeout: cfg.HttpWriteTimeout,
	}

	log.Info("starting http server",
		slog.String("address", cfg.HttpAddr),
	)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler

	switch env {
	case envLocal:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		break
	case envProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
		break
	}

	return slog.New(handler)
}
