package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/andreipimenov/golang-training-2021/internal/handler"
	"github.com/andreipimenov/golang-training-2021/internal/repository"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

func main() {
	r := chi.NewRouter()

	repo := repository.New()
	service := service.New(repo)
	h := handler.New(service)

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Method(http.MethodGet, "/price/{ticker}/{date}", h)
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		log.Println("Server is listening on :8080")
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-shutdown

	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped gracefully")
}
