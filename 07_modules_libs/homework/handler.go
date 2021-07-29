package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

func stockHandler(w http.ResponseWriter, r *http.Request) {

	token := "0WH6HZBAMK2FVZV2"

	ticker := chi.URLParam(r, "ticker")

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY&symbol=%s&apikey=%s", ticker, token)

	var highest, lowest, average, total float64

	response, err := http.Get(url)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := ioutil.ReadAll(response.Body)

	var s StockRequest
	err = json.Unmarshal(body, &s)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error)
		return
	}

	highest = s.Price[0].StockPrice
	lowest = s.Price[0].StockPrice

	for _, i := range s.Price {
		if i.StockPrice > highest {
			highest = i.StockPrice
		}
		if i.StockPrice < lowest {
			lowest = i.StockPrice
		}
		total += i.StockPrice
	}

	average = math.Floor((total/float64(len(s.Price)))*100) / 100 // round float to 2 decimal places for accuracy

	res := FormattedResponse{
		Ticker:        ticker,
		Highest_price: highest,
		Lowest_price:  lowest,
		Avg_price:     average,
	}

	writeResponse(w, http.StatusOK, res)
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	_, err := w.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
