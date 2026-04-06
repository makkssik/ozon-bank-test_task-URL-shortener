package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	redirectUseCase "url-shortener/internal/application/usecases/redirect"
	shortenUseCase "url-shortener/internal/application/usecases/shorten"
	"url-shortener/internal/config"
	memoRepo "url-shortener/internal/infrastructure/persistence/memory/repositories"
	pgRepo "url-shortener/internal/infrastructure/persistence/postgres/repositories"
	"url-shortener/internal/presentation/http/handlers/redirect"
	"url-shortener/internal/presentation/http/handlers/save"
	mwLogger "url-shortener/internal/presentation/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.MustLoad()

	storageType := flag.String("storage", "memory", "storage: memory / postgres")
	flag.Parse()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	var repository interface {
		shortenUseCase.URLSaver
		redirectUseCase.URLProvider
	}

	if *storageType == "postgres" {
		db, err := sql.Open("postgres", cfg.StoragePath)
		if err != nil {
			log.Error("failed to init postgres", slog.String("error", err.Error()))
			os.Exit(1)
		}

		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				log.Error("failed to close postgres", slog.String("error", err.Error()))
			} else {
				log.Info("postgres connection closed")
			}
		}(db)

		if err := db.Ping(); err != nil {
			log.Error("failed to ping postgres", slog.String("error", err.Error()))
			os.Exit(1)
		}

		repository = pgRepo.NewStorage(db)
	} else {
		repository = memoRepo.NewStorage()
	}

	shortenService := shortenUseCase.NewService(repository)
	redirectService := redirectUseCase.NewService(repository)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)

	router.Post("/", save.NewSaveHandler(log, shortenService))
	router.Get("/{alias}", redirect.NewRedirectHandler(log, redirectService))

	router.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger.json"),
	))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Info("server is running", slog.String("address", cfg.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

	sign := <-stop
	log.Info("stopping server", slog.String("signal", sign.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", slog.String("error", err.Error()))
	}

	log.Info("server stopped gracefully")
}
