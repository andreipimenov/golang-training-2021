package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func findTicket(w http.ResponseWriter, r *http.Request) {
	//get URL param from url and build GET request
	company := chi.URLParam(r, "company")
	date := chi.URLParam(r, "date")
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&outputsize=full&apikey=JC372XHXAL58Y0FL", company)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	//read data from url
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	//search data in json
	var i interface{}
	json.Unmarshal(body, &i)
	objects := i.(map[string]interface{})
	time := objects["Time Series (Daily)"].(map[string]interface{})
	Prices := time[date].(map[string]interface{})
	closePrice := Prices["4. close"].(string)
	result := fmt.Sprintf("{\"ticker\":\"%s\",\"close_price\":\"%s\",\"date\":\"%s\"}", company, closePrice, date)
	fmt.Println(result)
}

func main() {
	r := chi.NewRouter()
	r.MethodFunc(http.MethodGet, "/price/{company}/date/{date}/", findTicket)

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	srv.ListenAndServe()
}
