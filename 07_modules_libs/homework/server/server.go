package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rabbit72/golang-training-2021/07_modules_libs/homework/stock"
)

func getTickerStat(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")

	stockStat, err := stock.GetStockStat(ticker, apiKey)
	if err != nil {
		log.Printf("Cannot get stock data for %v - %v", ticker, err)
		http.Error(w, fmt.Errorf("the ticker %v has not been found", ticker).Error(), 404)
		return
	}
	payload, err := json.Marshal(stockStat)
	if err != nil {
		log.Printf("Unable to marshal ticker %v - %v", ticker, err)
		http.Error(w, err.Error(), 500)
	} else {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(payload))
		log.Printf("Successful response for %v ticker", ticker)
	}
}

func getAPIKey() string {
	envVarName := "STOCK_TOKEN"
	var apiKey string = os.Getenv(envVarName)
	if apiKey == "" {
		log.Fatalf("%v env variable is empty, fill the .env file before running", envVarName)
	}
	return apiKey
}

var apiKey string

func Serve(addr string) error {
	apiKey = getAPIKey()

	router := chi.NewRouter()
	router.Get("/price/{ticker}/stat", getTickerStat)
	err := http.ListenAndServe(addr, router)
	return err
}
