package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (c *aVClient) initAVClient(key string) {
	c.key = key
	c.client = &http.Client{
		Timeout: time.Duration(time.Minute * 5),
	}
}

// query to ALPHA VANTAGE
func (avClient aVClient) pricesByDate(ticker string, date string, ch chan result) {
	var priceResult result
	// sync by execution result
	defer func() {
		ch <- priceResult
	}()
	// parse url string with required date and ticker
	// this API get daily price info
	url_string := fmt.Sprintf("http://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%s&outputsize=full&apikey=%s", ticker, avClient.key)
	req, err := http.NewRequest(http.MethodGet, url_string, nil)
	if err != nil {
		log.Println(err)
		priceResult.err = err
		return
	}

	res, err := avClient.client.Do(req)
	if err != nil {
		log.Println(err)
		priceResult.err = err
		return
	}
	defer res.Body.Close()

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		priceResult.err = err
		return
	}

	if res.StatusCode != http.StatusOK {
		priceResult.err = errors.New("Bad status code")
		return
	}

	// parse returned data
	var jsonData interface{}
	err = json.Unmarshal(raw, &jsonData)
	if err != nil {
		priceResult.err = err
		return
	}
	dates, ok := jsonData.(map[string]interface{})["Time Series (Daily)"]
	if !ok {
		priceResult.err = errors.New("Cat't take time series")
		return
	}
	retDate, ok := dates.(map[string]interface{})[date]
	if !ok {
		priceResult.err = errors.New("There is no information on the specified date")
		return
	}

	prices := retDate.(map[string]interface{})
	priceResult.value.Ticker = ticker
	priceResult.value.Date = date
	priceResult.value.Open_price, ok = prices["1. open"].(string)
	if !ok {
		priceResult.value.Open_price = "Not found"
	}
	priceResult.value.High_price, ok = prices["2. high"].(string)
	if !ok {
		priceResult.value.High_price = "Not found"
	}
	priceResult.value.Low_price, ok = prices["3. low"].(string)
	if !ok {
		priceResult.value.Low_price = "Not found"
	}
	priceResult.value.Close_price, ok = prices["4. close"].(string)
	if !ok {
		priceResult.value.Close_price = "Not found"
	}
}
