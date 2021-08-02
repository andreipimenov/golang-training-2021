package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

var (
	errUnexpectedJSON = fmt.Errorf("unexpected json")
)

type stockAPIResponse map[time.Time]model.Price

func (s *stockAPIResponse) UnmarshalJSON(raw []byte) error {
	if s != nil && *s == nil {
		*s = make(map[time.Time]model.Price)
	}
	var i map[string]interface{}
	err := json.Unmarshal(raw, &i)
	if err != nil {
		return err
	}

	//if there error from API
	errAPI, ok := i["Error Message"].(string)
	if ok {
		return fmt.Errorf("error from API: %v", errAPI)
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
		(*s)[d] = model.Price{Open: open, High: high, Low: low, Close: close}
	}
	return nil
}
