package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func GetStockInfo(w http.ResponseWriter, r *http.Request) {

	symbol := chi.URLParam(r, "ticker")
	token := "99NO7A1MXMA2O24S"
	//token  = "PLKTKS28UM23OXGY"
	requestHiLo := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=%s&interval=60min&apikey=%s", symbol, token)
	requestAvg := fmt.Sprintf("https://www.alphavantage.co/query?function=SMA&symbol=%s&interval=daily&time_period=10&series_type=open&apikey=%s", symbol, token)
	intradayInfo, errIntraday := http.Get(requestHiLo)
	if errIntraday != nil {
		writeResponse(w, http.StatusBadRequest, Error{errIntraday.Error()})
		return
	}
	avgInfo, errAvg := http.Get(requestAvg)
	if errAvg != nil {
		writeResponse(w, http.StatusBadRequest, Error{errAvg.Error()})
		return
	}

	data, errReadAll := ioutil.ReadAll(intradayInfo.Body)
	if errReadAll != nil {
		writeResponse(w, http.StatusBadRequest, Error{errReadAll.Error()})
		return
	}

	dataSMA, errReadAllSMA := ioutil.ReadAll(avgInfo.Body)
	if errReadAllSMA != nil {
		writeResponse(w, http.StatusBadRequest, Error{errReadAllSMA.Error()})
		return
	}

	var unmData RequestJSONIntraDay
	errUnmarshal := json.Unmarshal(data, &unmData)

	if errUnmarshal != nil {
		writeResponse(w, http.StatusBadRequest, Error{errUnmarshal.Error()})
		return
	}

	if unmData.Note != "" {
		writeResponse(w, http.StatusBadRequest, Error{unmData.Note})
		return
	}

	if unmData.M.TimeZone == "" {
		writeResponse(w, http.StatusBadRequest, Error{"Wrong ticker name"})
		return
	}

	currenDate := unmData.M.LastRefreshed
	if len(currenDate) < 10 {
		writeResponse(w, http.StatusBadRequest, Error{"wrong date value"})
		return
	}

	globalMax := 0.0
	globalMin := math.MaxFloat64
	currenDate = currenDate[0:10]
	for k, v := range unmData.TS {
		if strings.Contains(k, currenDate) {
			localMax, errParseMax := strconv.ParseFloat(v.High, 64)
			if errParseMax != nil {
				writeResponse(w, http.StatusBadRequest, errParseMax)
				return
			}
			localMin, errParseMin := strconv.ParseFloat(v.Low, 64)
			if errParseMin != nil {
				writeResponse(w, http.StatusBadRequest, errParseMin)
				return
			}
			if globalMin > localMin {
				globalMin = localMin
			}
			if globalMax < localMax {
				globalMax = localMax
			}
		}
	}

	var unmSMA RequestJSONIndicator
	errUnmarshal = json.Unmarshal(dataSMA, &unmSMA)

	if unmSMA.Note != "" {
		writeResponse(w, http.StatusBadRequest, Error{unmSMA.Note})
		return
	}

	if errUnmarshal != nil {
		writeResponse(w, http.StatusBadRequest, Error{errUnmarshal.Error()})
		return
	}

	sma, ok := unmSMA.Data[currenDate]
	if !ok {
		writeResponse(w, http.StatusBadRequest, Error{"cant find current date avg"})
		return
	}

	avg, errParse := strconv.ParseFloat(sma.SMA, 64)
	if errParse != nil {
		writeResponse(w, http.StatusBadRequest, errParse)
		return
	}

	result := NewStock(symbol, avg, globalMax, globalMin)

	err := avgInfo.Body.Close()
	if err != nil {
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}
	err = intradayInfo.Body.Close()
	if err != nil {
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}

	writeResponse(w, http.StatusOK, result)

}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write([]byte(b))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
