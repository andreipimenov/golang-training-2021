package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
)

func stockHandler(w http.ResponseWriter, r *http.Request) {

	token := "0WH6HZBAMK2FVZV2"

	ticker := chi.URLParam(r, "ticker")

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_MONTHLY&symbol=%s&apikey=%s", ticker, token)

	response, err := http.Get(url)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err2 := ioutil.ReadAll(response.Body)

	var s interface{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err2 != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, s)
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write(b)
}
