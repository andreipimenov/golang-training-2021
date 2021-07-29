package main

import (
	"encoding/json"
	"log"
	"strconv"
)

type tickerInfo struct {
	StockPrice float64
}

type StockRequest struct {
	Price []tickerInfo
}

type FormattedResponse struct {
	Ticker        string
	Highest_price float64
	Lowest_price  float64
	Avg_price     float64
}

func (sr *StockRequest) UnmarshalJSON(b []byte) error {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal(err)
	}

	rawResponse := f.(map[string]interface{})

	data := rawResponse["Monthly Time Series"].(map[string]interface{})

	var openPrice []tickerInfo
	var currentDayPrice tickerInfo
	var parsedData map[string]interface{}

	for _, d := range data {
		parsedData = d.(map[string]interface{})

		currentDayPrice = tickerInfo{}
		// we use "open" price only in demo purposes. It`s not accurate data for calculating and highest, lowest values
		// should parse high and low values separately then calculate average price with ((high + low)/2)
		currentDayPrice.StockPrice, _ = strconv.ParseFloat(parsedData["1. open"].(string), 64)

		openPrice = append(openPrice, currentDayPrice)
	}

	sr.Price = openPrice

	return nil
}

