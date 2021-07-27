package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Stock struct {
	Pagination struct {
	} `json:"pagination"`
	Data []struct {
		Open        float64     `json:"open"`
		High        float64     `json:"high"`
		Low         float64     `json:"low"`
		Close       float64     `json:"close"`
		Volume      float64     `json:"volume"`
		AdjHigh     interface{} `json:"adj_high"`
		AdjLow      interface{} `json:"adj_low"`
		AdjClose    float64     `json:"adj_close"`
		AdjOpen     interface{} `json:"adj_open"`
		AdjVolume   interface{} `json:"adj_volume"`
		SplitFactor float64     `json:"split_factor"`
		Symbol      string      `json:"symbol"`
		Exchange    string      `json:"exchange"`
		Date        string      `json:"date"`
	} `json:"data"`
}

func getClosePrice(ticker string, date string) (float64, error) {
	var stock Stock
	resp, err := http.Get("http://api.marketstack.com/v1/eod?access_key=3b795be41caab45049f693481cce6df4&symbols=" +
		ticker + "&date_from=" +
		date + "&date_to=" +
		date)
	body, err := ioutil.ReadAll(resp.Body)
	if errorHandler(err) {
		return 0, err
	}
	err = json.Unmarshal(body, &stock)
	if errorHandler(err) || len(stock.Data) == 0 {
		if len(stock.Data) == 0 {
			err = errors.New("error in get request to api.marketstack")
		}
		return 0, err
	}
	return stock.Data[0].Close, nil
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write(b)
}

func errorHandler(err error) bool {
	if err != nil {
		return true
	} else {
		return false
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	firstDate := chi.URLParam(r, "first_date")
	secondDate := chi.URLParam(r, "second_date")
	dayOne, err := getClosePrice(ticker, firstDate)
	if errorHandler(err) {
		writeResponse(w, http.StatusBadRequest, "error occurred in request to api.marketstack")
		return
	}
	dayTwo, err := getClosePrice(ticker, secondDate)
	if errorHandler(err) {
		writeResponse(w, http.StatusBadRequest, "error occurred in request to api.marketstack")
		return
	}
	diff := dayTwo - dayOne
	percentageDiff := diff / dayOne * 100
	diffMap := map[string]float64{"diffPercantage": percentageDiff}
	writeResponse(w, http.StatusOK, diffMap)
}

func main() {
	router := chi.NewRouter()
	router.MethodFunc(http.MethodGet, "/price/{ticker}/diff/{first_date}/{second_date}", get)
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		log.Println("Server is listening on :8080")
		err := server.ListenAndServe()
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
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped gracefully")
}
