package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//Usage: GET /price/{ticker}/stat
func stockHandler(w http.ResponseWriter, r *http.Request) {

	ticker := chi.URLParam(r, "ticker")
	//API key in code only for demo purposes
	//Will be invalidated 01/08/21
	api_token := "KCVH80JGUFNYE8L9"

	url_string := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY_ADJUSTED&symbol=%s&apikey=%s", ticker, api_token)

	resp, err := http.Get(url_string)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}

	var s StockRequest
	err = json.Unmarshal(body, &s)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}

	var (
		highest float64
		lowest  float64
		average float64
		total   float64
	)

	highest = s.StockData[0].Close
	lowest = s.StockData[0].Close
	total = 0

	//Loop for getting the highest, lowest and total values
	for _, v := range s.StockData {
		total += v.Close
		if v.Close > highest {
			highest = v.Close
		}
		if v.Close < lowest {
			lowest = v.Close
		}
	}

	//Rounding trick to truncate output
	average = math.Round((total/float64(len(s.StockData)))*10) / 10

	fmt.Printf("Highest: %f, Lowest: %f, Average: %f\n", highest, lowest, average)

	res := StockResponse{
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
	w.Write([]byte(b))
}
