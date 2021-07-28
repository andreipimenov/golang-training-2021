package stock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type StockStat struct {
	Ticker       string  `json:"ticker"`
	HighestPrice float64 `json:"highest_price,string"`
	LowestPrice  float64 `json:"lowest_price,string"`
	AvgPrice     float64 `json:"avg_price,string"`
}

func GetStockStat(stockTicker string, apiKey string) (*StockStat, error) {
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=OVERVIEW&symbol=%v&apikey=%v",
		strings.ToUpper(stockTicker), apiKey,
	)
	response, err := http.Get(url)
	if err != nil {
		return &StockStat{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorMessage, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return &StockStat{}, err
		}
		return &StockStat{}, fmt.Errorf("%v", errorMessage)

	}

	rawJSON, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &StockStat{}, err
	}

	var stockStat StockStat
	if err := json.Unmarshal(rawJSON, &stockStat); err != nil {
		return &StockStat{}, err
	}

	return &stockStat, nil
}

func (stock *StockStat) UnmarshalJSON(b []byte) error {
	const (
		symbolField   = "Symbol"
		weekHighField = "52WeekHigh"
		weekLowField  = "52WeekLow"
		avgField      = "200DayMovingAverage"
	)
	mapStockStat := make(map[string]string)
	if err := json.Unmarshal(b, &mapStockStat); err != nil {
		return err
	}

	if len(mapStockStat) == 0 {
		return fmt.Errorf("JSON contains zero keys")
	}

	ticker, ok := mapStockStat[symbolField]
	if !ok {
		return fmt.Errorf("%v key not found in JSON", symbolField)
	}

	highestPrice, err := exctractAsFloat(mapStockStat, weekHighField)
	if err != nil {
		return err
	}

	lowestPrice, err := exctractAsFloat(mapStockStat, weekLowField)
	if err != nil {
		return err
	}

	avgPrice, err := exctractAsFloat(mapStockStat, avgField)
	if err != nil {
		return err
	}

	stock.Ticker = ticker
	stock.HighestPrice = highestPrice
	stock.LowestPrice = lowestPrice
	stock.AvgPrice = avgPrice
	return nil
}

func exctractAsFloat(m map[string]string, key string) (float64, error) {
	plainValue, ok := m[key]
	if !ok {
		return 0., fmt.Errorf("%v key not found in JSON", key)
	}
	value, err := strconv.ParseFloat(plainValue, 64)
	if err != nil {
		return 0., fmt.Errorf("%v cannot be parsed as float64", key)
	}
	return value, nil
}
