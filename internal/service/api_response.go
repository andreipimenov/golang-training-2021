package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

var (
	errUnexpectedJSON   = fmt.Errorf("unexpected json")
	errUnexpectedDate   = fmt.Errorf("unexprected date")
	errWeekendDate      = fmt.Errorf("weekend date")
	errDateDoesNotExist = fmt.Errorf("date does not exist")
)

type StockAPIResponse struct {
	TimeSeriesDaily   map[time.Time]model.Price
	LastRefreshedTime time.Time
}

func (s *StockAPIResponse) UnmarshalJSON(raw []byte) error {
	var i map[string]interface{}
	err := json.Unmarshal(raw, &i)
	if err != nil {
		return err
	}

	lastRefreshed, ok := i["Meta Data"].(map[string]interface{})["3. Last Refreshed"].(string)
	if !ok {
		return errUnexpectedJSON
	}

	lastRefreshedTime, err := time.Parse("2006-01-02", lastRefreshed)
	if err != nil {
		return err
	}

	tsd, ok := i["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		return errUnexpectedJSON
	}

	for k, v := range tsd {
		d, err := time.Parse("2006-01-02", k)
		if err != nil {
			return err
		}
		x, ok := v.(map[string]interface{})
		if !ok {
			return errUnexpectedJSON
		}
		open, ok := x["1. open"].(string)
		if !ok {
			return errUnexpectedJSON
		}
		high, ok := x["2. high"].(string)
		if !ok {
			return errUnexpectedJSON
		}
		low, ok := x["3. low"].(string)
		if !ok {
			return errUnexpectedJSON
		}
		close, ok := x["4. close"].(string)
		if !ok {
			return errUnexpectedJSON
		}
		s.TimeSeriesDaily[d] = model.Price{Open: open, High: high, Low: low, Close: close}
	}
	s.LastRefreshedTime = lastRefreshedTime

	return nil
}
