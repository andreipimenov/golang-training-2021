// homework implements http server with one endpoint
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Alpha Vantage API Key
// YFG2RTQZ5KYRAUP4
const (
	APIKEY = "YFG2RTQZ5KYRAUP4"
	URL    = "https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY_ADJUSTED&outputsize=full"
)

type Error struct {
	Error string
}

type Model struct {
	Ticker     string `json:"ticker"`
	ClosePrice string `json:"close_price"`
	Date       string `json:"date"`
}

// V1 returns the stock price for exact ticker on exact date
// GET /price/AAPL/date/2021-07-2
// GET /price/{ticker}/stat
func handleStockPrice(w http.ResponseWriter, r *http.Request) {
	// Getting params
	ticker := chi.URLParam(r, "ticker")
	date := chi.URLParam(r, "date")
	url := fmt.Sprintf("%s&symbol=%s&apikey=%s", URL, ticker, APIKEY)
	// url := "https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY_ADJUSTED&outputsize=full&symbol=IBM&apikey=YFG2RTQZ5KYRAUP4"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// Now we have 2 read data from url
	// Base64 format
	bodyInBase64, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var bodyParsed interface{}
	err = json.Unmarshal(bodyInBase64, &bodyParsed)
	if err != nil {
		log.Println(err)
		writeRsponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}

	objects, ok := bodyParsed.(map[string]interface{})
	if !ok {
		writeRsponse(w, http.StatusNotFound, Error{"404_1"})
		return
	}
	// writeRsponse(w, http.StatusOK, objects)
	firstObject, ok := objects["Monthly Adjusted Time Series"].(map[string]interface{})
	if !ok {
		writeRsponse(w, http.StatusNotFound, Error{"404_2"})
		return
	}

	// writeRsponse(w, http.StatusOK, firstObject)

	certainDate, ok := firstObject[date].(map[string]interface{})
	if !ok {
		writeRsponse(w, http.StatusNotFound, Error{"404_3 's no such date price"})
		return
	}
	// writeRsponse(w, http.StatusOK, certainDate["4. close"])
	closePrice, ok := certainDate["4. close"].(string)
	if !ok {
		writeRsponse(w, http.StatusNotFound, Error{"There's no close price"})
		return
	}
	// writeRsponse(w, http.StatusOK, closePrice)
	///*
	//	Return answer in
	//		{
	//			"ticker":"AAPL",
	//			"close_price":"146.80",
	//			"date":"2021-07-22"
	//		}
	//	format from valid JSON
	//*/
	company := Model{
		Ticker:     ticker,
		ClosePrice: closePrice,
		Date:       date,
	}
	writeRsponse(w, http.StatusOK, company)
}

// V2 returns the highest, lowest and average prices
// GET /price/AAPL/stat
func handleHALPrices(w http.ResponseWriter, r *http.Request) {
	// TODO : Implement
}

// V3 returns close prices percentage difference
func handlePercentageDifference(w http.ResponseWriter, r *http.Request) {
	// TODO : Implement
}

// writeRsponse - is a helper func to get data
func writeRsponse(w http.ResponseWriter, statusCode int, value interface{}) {
	b, _ := json.Marshal(value)
	w.WriteHeader(statusCode)
	//Send data through network
	w.Write([]byte(b))
}

func main() {
	// Create New Router
	router := chi.NewRouter()

	// Map endpoint & event handler
	router.Get("/price/{ticker}/date/{date}", handleStockPrice)

	// Declare http server
	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Gracefull shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	// Clearing of resources
	defer signal.Stop(shutdown)

	go func() {
		log.Println("Server is listening on :8080")
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	// Wait for Interrupt signal
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
