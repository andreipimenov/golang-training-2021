package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

var (
	errUnexpectedJSON = fmt.Errorf("unexpected json")
)

type stockAPIResponse map[time.Time]model.Price

type dataModel struct {
	Open  string `json:"1. open"`
	High  string `json:"2. high"`
	Low   string `json:"3. low"`
	Close string `json:"4. close"`
}

func (s *stockAPIResponse) UnmarshalJSON(raw []byte) error {
	if s != nil && *s == nil {
		*s = make(map[time.Time]model.Price)
	}
	var i map[string]interface{}
	err := json.Unmarshal(raw, &i)
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
			log.Println(err)
			return err
		}

		body, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
			return err
		}

		data := dataModel{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Println(err)
			return errUnexpectedJSON
		}

		(*s)[d] = model.Price{Open: data.Open, High: data.High, Low: data.Low, Close: data.Close}
	}
	return nil
}
