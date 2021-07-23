package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Error struct {
	Error string
}

type StockRequest struct {
	StockData []StockData
}

type StockData struct {
	Close float64
}

type StockResponse struct {
	Ticker        string  `json:"ticker"`
	Highest_price float64 `json:"highest_price"`
	Lowest_price  float64 `json:"lowest_price"`
	Avg_price     float64 `json:"avg_price"`
}

//Custom unmarshal to properly unmarshal AlphaVantage nested JSON
//and get the data we need
func (sr *StockRequest) UnmarshalJSON(b []byte) error {

	//Unmarshal to an interface because the number of nested objects depends on the ticker specified in the API call
	//and therefore the structure of the response is unknown in advance of the call.
	var f interface{}
	json.Unmarshal(b, &f)
	root_map := f.(map[string]interface{})

	//If an unknown tickers info is requested through AlphaVantage API it responds with
	//a single "Error Message" property. Here I check for this and return with an error if it's the case.
	if len(root_map) == 1 {
		return fmt.Errorf("unknown ticker")
	}

	//Iteration over nested map[string]interface{} to get the required data.
	mats_map := root_map["Monthly Adjusted Time Series"].(map[string]interface{})
	var stock_arr []StockData
	var temp StockData
	var res_map map[string]interface{}
	for _, v := range mats_map {

		res_map = v.(map[string]interface{})
		temp = StockData{}
		temp.Close, _ = strconv.ParseFloat(res_map["5. adjusted close"].(string), 64)

		stock_arr = append(stock_arr, temp)
	}

	sr.StockData = stock_arr

	return nil
}
