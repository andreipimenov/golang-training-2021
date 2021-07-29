package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func convert(x string) float64 {
	res, err := strconv.ParseFloat(strings.Trim(x, " "), 64)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
func countPercentage(x, y string) float64 {
	fst, sec := convert(x), convert(y)
	return math.Round((sec-fst)/fst*10000) / 100
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(code)
	w.Write([]byte(b))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	key := "QPYBY2JSO2PMHOCX"
	ticker := chi.URLParam(r, "ticker")
	firstDate := chi.URLParam(r, "first_date")
	secondDate := chi.URLParam(r, "second_date")
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%v&apikey=%v", ticker, key)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	connectField := &connectField{}

	err = json.Unmarshal(raw, &connectField)

	if err != nil {
		log.Fatal(err)
	}

	diff := countPercentage(connectField.DailyRes[firstDate].Price, connectField.DailyRes[secondDate].Price)

	respField := respStr{
		Ticker:         ticker,
		PercentageDiff: diff,
		FirstDate:      firstDate,
		SecondDate:     secondDate,
	}

	writeResponse(w, http.StatusOK, respField)
}
