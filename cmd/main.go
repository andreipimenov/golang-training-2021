package main

import (
	"context"
	"fmt"
	cfg "github.com/andreipimenov/golang-training-2021/internal/config"
	"github.com/andreipimenov/golang-training-2021/internal/handler"
	"github.com/andreipimenov/golang-training-2021/internal/repository"
	"github.com/andreipimenov/golang-training-2021/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	loggerInit()
	repo, clz := repository.Get()
	defer clz()
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.RequestLogger(&handler.LogFormatter{Logger: &log.Logger}))
		r.Use(middleware.Recoverer)
		r.Method(http.MethodGet,
			handler.Path,
			handler.New(&log.Logger, service.New(&log.Logger, repo, cfg.Get().ExternalAPIToken)))
	})
	serverRun(r)
}

func loggerInit() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.Level(cfg.Get().LogLevel))
}

func serverRun(r *chi.Mux) {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Get().Port),
		Handler: r,
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)
	go func() {
		log.Info().Msgf("Server is listening on :%d", cfg.Get().Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server error")
		}
	}()
	<-shutdown
	log.Info().Msg("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server shutdown error")
	}
	log.Info().Msg("Server stopped gracefully")
}
