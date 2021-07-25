package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// result of successful query
type jsonRes struct {
	Ticker      string      `json:"ticker"`
	Date        string      `json:"date"`
	Open_price  interface{} `json:"open_price"`
	High_price  interface{} `json:"high_price"`
	Low_price   interface{} `json:"low_price"`
	Close_price interface{} `json:"close_price"`
}

// result of the query
// if query is successful result will be placed in "value"
// or you will get error
type result struct {
	value jsonRes
	err   error
}

type customError struct {
	Error string
}

//handle func for endpoint /price/{ticker}/date/{date}
func (s StockServer) price(w http.ResponseWriter, r *http.Request) {

	ticker := chi.URLParam(r, "ticker")
	date := chi.URLParam(r, "date")

	ctx := r.Context()
	readyCh := make(chan result)

	go s.avClient.pricesByDate(ticker, date, readyCh)

	select {
	case answer := <-readyCh:
		if answer.err != nil {
			writeResponse(w, http.StatusBadRequest, customError{answer.err.Error()})
			return
		}
		writeResponse(w, http.StatusOK, answer.value)
	case <-ctx.Done():
		log.Printf("Interrupted with err: %v", ctx.Err())
	}
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write([]byte(b))
}

//run StockServer
func (server *StockServer) Run(port interface{}) {

	var strPort string
	switch port.(type) {
	case int:
		strPort = strconv.Itoa(port.(int))
	case string:
		strPort = port.(string)
	default:
		log.Fatal("Bad port")
	}

	server.avClient.initAVClient(apiKey)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/price/{ticker}/date/{date}", server.price)

	server.httpServer = http.Server{
		Addr:    ":" + strPort,
		Handler: r,
	}

	go func() {
		log.Println("Server is listening on :8080")
		err := server.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

//wait StockServer shutdown
func (server StockServer) WaitShutdown() {

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer func() {
		cancel()
	}()
	if err := server.httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped gracefully")
}
