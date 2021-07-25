package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write([]byte(b))
}

func validateDiffDates(first, second string) error {
	// check if dates are valid
	f, err := time.Parse(dateLayoutISO, first)
	if err != nil {
		return fmt.Errorf("cannot parse first_date parameter: %v", err)
	}
	s, err := time.Parse(dateLayoutISO, second)
	if err != nil {
		return fmt.Errorf("cannot parse second_date parameter: %v", err)
	}
	// check if firstDate < secondDate
	if !f.Before(s) {
		return fmt.Errorf("first_date should be less than second_date")
	}
	return nil
}

func GetDiffAsync(chanErr chan (error), td *TickerDiff) {

	// Usage: https://www.alphavantage.co/documentation/#intraday-extended
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%v&outputsize=%v&apikey=%v",
		td.Ticker, td.Format, apiKey)
	// Making  http-request
	httpClient := &http.Client{
		Timeout: time.Duration(time.Second * 15),
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		chanErr <- err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		chanErr <- fmt.Errorf("alphavantage call error: %v", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		chanErr <- fmt.Errorf("alphavantage call error (not-OK)")
	}

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		chanErr <- fmt.Errorf("alphavantage call error: %v", err.Error())
	}
	// Parsing data
	tsd := TimeSeriesDaily{}
	err = json.Unmarshal(raw, &tsd)
	if err != nil {
		chanErr <- fmt.Errorf("alphavantage response unmarshal error: %v", err.Error())
	}
	if tsd.ErrorMessage != "" {
		chanErr <- fmt.Errorf("alphavantage response unmarshal error: %v", tsd.ErrorMessage)
	}
	// Let's check if first and second dates are in response and change them if not
	// fmt.Println(td.FirstDate, td.SecondDate)
	newFirstDate := checkDateInAVResp(td.FirstDate, &tsd)
	newSecondDate := checkDateInAVResp(td.SecondDate, &tsd)
	// fmt.Println(newFirstDate, newSecondDate)
	// and calculating perc diff
	f, _ := strconv.ParseFloat(tsd.TimeSeriesData[newFirstDate].Close, 64)
	s, _ := strconv.ParseFloat(tsd.TimeSeriesData[newSecondDate].Close, 64)
	td.Diff = fmt.Sprintf("%.2f", ((s-f)/f)*100)
	chanErr <- nil
}

// checkDateInAVResp allows to check if a date exists in the AV response
// it would be missed if it's weekends day or something like that
// If we don't have this date, we'll dicrease it to the closest existing date
// So, we'll return a usable date in any case
func checkDateInAVResp(date string, tsd *TimeSeriesDaily) string {
	if _, ok := tsd.TimeSeriesData[date]; !ok {
		d, _ := time.Parse(dateLayoutISO, date)
		d = d.Add(-24 * time.Hour)
		return checkDateInAVResp(d.Format(dateLayoutISO), tsd)
	} else {
		return date
	}
}

// func GetDiff(td *TickerDiff) error {

// 	// Usage: https://www.alphavantage.co/documentation/#intraday-extended
// 	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%v&outputsize=%v&apikey=%v",
// 		td.Ticker, td.Format, apiKey)
// 	// Making  http-request
// 	httpClient := &http.Client{
// 		Timeout: time.Duration(time.Second * 15),
// 	}
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return err
// 	}

// 	res, err := httpClient.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("alphavantage call error: %v", err.Error())
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return fmt.Errorf("alphavantage call error (not-OK)")
// 	}

// 	raw, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return fmt.Errorf("alphavantage call error: %v", err.Error())
// 	}
// 	// Parsing data
// 	tsd := TimeSeriesDaily{}
// 	err = json.Unmarshal(raw, &tsd)
// 	if err != nil {
// 		return fmt.Errorf("alphavantage response unmarshal error: %v", err.Error())
// 	}
// 	if tsd.ErrorMessage != "" {
// 		return fmt.Errorf("alphavantage response unmarshal error: %v", tsd.ErrorMessage)
// 	}
// 	// Let's check if first and second dates are in response and change them if not
// 	// fmt.Println(td.FirstDate, td.SecondDate)
// 	newFirstDate := checkDateInAVResp(td.FirstDate, &tsd)
// 	newSecondDate := checkDateInAVResp(td.SecondDate, &tsd)
// 	// fmt.Println(newFirstDate, newSecondDate)
// 	// and calculating perc diff
// 	f, _ := strconv.ParseFloat(tsd.TimeSeriesData[newFirstDate].Close, 64)
// 	s, _ := strconv.ParseFloat(tsd.TimeSeriesData[newSecondDate].Close, 64)
// 	td.Diff = fmt.Sprintf("%.2f", ((s-f)/f)*100)
// 	return nil
// }
