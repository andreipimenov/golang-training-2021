package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type HandleReq struct {
	Data interface{}
}

func (h *HandleReq) Get(w http.ResponseWriter, r *http.Request) {
	const API_KEY = "HXPQ9N2CBVG2W6MJ"
	ticker := chi.URLParam(r, "ticker")
	date := chi.URLParam(r, "date")
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&apikey=%s", ticker, API_KEY)

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(raw, &h.Data)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, Error{"Not found!"})
		return
	}

	dataAll := h.Data.(map[string]interface{})
	metaData := dataAll["Meta Data"].(map[string]interface{})
	ticker = metaData["2. Symbol"].(string)
	dataTime := dataAll["Time Series (Daily)"].(map[string]interface{})
	adjClosePrice := dataTime[date].(map[string]interface{})["4. close"]
	res, _ := strconv.ParseFloat(adjClosePrice.(string), 64)

	info := RequestInfo{
		Ticker: ticker,
		Price:  res,
		Date:   date,
	}

	fmt.Printf("{\"ticker\": \"%v\", \"close_price\": \"%.2f\", \"date\": %v}\n", info.Ticker, info.Price, info.Date)
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write([]byte(b))
}
