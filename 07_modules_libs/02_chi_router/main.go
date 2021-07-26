package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	userHandler := &UserHandler{data: sync.Map{}}

	r.Get("/ping", ping)
	r.Get("/long", bigSyncJob)
	r.MethodFunc(http.MethodPost, "/user", userHandler.Post)
	r.MethodFunc(http.MethodGet, "/user/{userID}", userHandler.Get)

	// r.Route("/", func(r chi.Router) {
	//   r.Use(middleware.Logger)
	// })

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
