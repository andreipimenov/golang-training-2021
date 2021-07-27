package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"time"
)

const ApiKey = "2NVMJS308BND1F0Y"

func getStockInfo(w http.ResponseWriter, r *http.Request) {

	stockName := chi.URLParam(r, "Ticker")

	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=%s&apikey=%s", stockName, ApiKey)

	response, error := http.Get(url)

	if error != nil {
		writeResponse(w, http.StatusBadRequest, error.Error())
		return
	}

	body, error2 := ioutil.ReadAll(response.Body)

	var s interface{}
	error = json.Unmarshal(body, &s)

	if error != nil || error2 != nil {
		writeResponse(w, http.StatusBadRequest, error.Error())
		return
	}

	objects := s.(map[string]interface{})

	timeSeriesDaily := objects["Time Series (Daily)"].(map[string]interface{})

	countDays := 1
	for {
		date := time.Now().AddDate(0, 0, -countDays).Format("2006-01-02")
		countDays++

		if dayData, accesError := timeSeriesDaily[date]; accesError {

			data := dayData.(map[string]interface{})

			type info struct {
				Ticker        string `json:"stockName"`
				Highest_price string `json:"highest_price"`
				Lowest_price  string `json:"lowest_price"`
				Avg_price     string `json:"avg_price"`
			}

			tickerInfo := info{
				Ticker:        stockName,
				Highest_price: data["2. high"].(string),
				Lowest_price:  data["3. low"].(string),
				Avg_price:     data["5. adjusted close"].(string),
			}

			fmt.Println(tickerInfo)

			writeResponse(w, http.StatusOK, tickerInfo)
			return
		}

		if countDays == len(timeSeriesDaily)-1 {
			writeResponse(w, http.StatusBadRequest, error.Error())
			return
		}
	}

}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write(b)
}
