package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/andreipimenov/golang-training-2021/internal/config"
	"github.com/andreipimenov/golang-training-2021/internal/handler"
	"github.com/andreipimenov/golang-training-2021/internal/repository"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Configuration error")
	}

	r := chi.NewRouter()

	db, err := sql.Open("postgres", cfg.DBConnString)
	if err != nil {
		logger.Fatal().Err(err).Msg("DB initializing error")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal().Err(err).Msg("DB pinging error")
	}

	stockRepo := repository.NewCache()
	stockService := service.NewStock(&logger, stockRepo, cfg.ExternalAPIToken)
	stockHandler := handler.NewStock(&logger, stockService)

	authRepo := repository.NewAuth(db)
	authService := service.NewAuth(&logger, authRepo, []byte(cfg.Secret))
	authHandler := handler.NewAuth(&logger, authService)

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.RequestLogger(&handler.LogFormatter{Logger: &logger}))
		r.Use(middleware.Recoverer)
		r.Use(handler.JWT([]byte(cfg.Secret)))
		r.Method(http.MethodGet, handler.StockPath, stockHandler)
		r.Method(http.MethodPost, handler.AuthPath, authHandler)
	})

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		logger.Info().Msgf("Server is listening on :%d", cfg.Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server error")
		}
	}()

	<-shutdown

	logger.Info().Msg("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server shutdown error")
	}

	logger.Info().Msg("Server stopped gracefully")
}
