package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type AVantageCacheHandler struct {
	data sync.Map
}

type AVantage struct {
	MetaData        map[string]string      `json:"Meta Data"`
	TimeSeriesDaily map[string]interface{} `json:"Time Series (Daily)"`
}

type AVantageResp struct {
	Ticker     string `json:"ticker"`
	PercDiff   string `json:"percentage_diff"`
	FirstData  string `json:"first_data"`
	SecondData string `json:"second_data"`
}

type Error struct {
	Error string
}

func (c *AVantageCacheHandler) Get(w http.ResponseWriter, r *http.Request) {

	apiKey := "ODFNBOJQFCDO91TF"
	ticker := chi.URLParam(r, "ticker")
	firstDate := chi.URLParam(r, "first_date")
	secondDate := chi.URLParam(r, "second_date")

	if _, ok := c.data.Load(ticker); ok {
		log.Println("from cache")
		writeResponse(w, http.StatusOK, resp(c, firstDate, secondDate, ticker))
		return
	}

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", ticker, apiKey)

	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode == http.StatusOK {
		aVantage := &AVantage{}
		err := json.Unmarshal(raw, &aVantage)
		if err != nil {
			log.Fatal(err)
		}
		c.data.Store(ticker, aVantage.TimeSeriesDaily)
		writeResponse(w, http.StatusOK, resp(c, firstDate, secondDate, ticker))
	}

}

func resp(c *AVantageCacheHandler, firstDate, secondDate, ticker string) *AVantageResp {

	resp := &AVantageResp{}

	m, _ := c.data.Load(ticker)
	t, _ := m.(map[string]interface{})

	fd, _ := t[firstDate].(map[string]interface{})["4. close"].(string)
	firstPrice, _ := strconv.ParseFloat(fd, 64)

	sd, _ := t[secondDate].(map[string]interface{})["4. close"].(string)
	secondPrice, _ := strconv.ParseFloat(sd, 64)

	diff := (secondPrice - firstPrice) / firstPrice * 100

	resp.PercDiff = strconv.FormatFloat(diff, 'f', 2, 64)
	resp.FirstData = firstDate
	resp.SecondData = secondDate
	resp.Ticker = ticker

	return resp
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(b))
}
