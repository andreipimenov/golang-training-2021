package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//Testing
//http://127.0.0.1:8080/price/AAPL/stat
//http://127.0.0.1:8080/price/IBM/stat
//http://127.0.0.1:8080/price/GOOGL/stat

func main() {
	router := chi.NewRouter()

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.MethodFunc(http.MethodGet, "/price/{Ticker}/stat", getStockInfo)

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
