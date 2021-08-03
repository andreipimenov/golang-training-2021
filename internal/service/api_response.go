package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

var (
	errUnexpectedJSON = fmt.Errorf("unexpected json")
	errInvalidParam   = fmt.Errorf("invalid parameters")
)

type stockAPIResponse model.Ticker

func (s *stockAPIResponse) UnmarshalJSON(raw []byte) error {
	var i map[string]interface{}
	err := json.Unmarshal(raw, &i)
	if err != nil {
		return err
	}
	_, ok := i["Error Message"]
	if ok {
		return errInvalidParam
	}
	metadata, ok := i["Meta Data"].(map[string]interface{})
	if !ok {
		return errUnexpectedJSON
	}

	s.Name, ok = metadata["2. Symbol"].(string)
	if !ok {
		return errUnexpectedJSON
	}

	lastRefreshed, ok := metadata["3. Last Refreshed"].(string)
	if !ok {
		return errUnexpectedJSON
	}

	s.LastRefreshed, err = time.Parse("2006-01-02", lastRefreshed)
	if err != nil {
		return err
	}

	tsd, ok := i["Time Series (Daily)"].(map[string]interface{})
	if !ok {
		return errUnexpectedJSON
	}

	history := make(map[time.Time]model.Price)
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
		history[d] = model.Price{Open: open, High: high, Low: low, Close: close}
	}
	s.History = history
	return nil
}
