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
)

// Implement http server with one endpoint GET /price/{ticker}/stat
// which returns the highest, lowest and average prices
//(e.g. GET /price/AAPL/stat -> {"ticker":"AAPL","highest_price":"149.80","lowest_price":"34.80","avg_price":"97.56"})

// To obtain stock price info you can use external API (e.g. https://www.alphavantage.co/)

func main() {

	r := chi.NewRouter()

	r.Get("/price/{ticker}/stat", stockHandler)

	port := ":8080"

	srv := http.Server{
		Addr:    port,
		Handler: r,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		log.Println("Server is listening on ", port)
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
