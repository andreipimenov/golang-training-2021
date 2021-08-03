package alphavantage

import (
	"encoding"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type Prices map[time.Time]model.Price

type date time.Time

func (d *date) UnmarshalText(raw []byte) error {
	t, err := time.Parse("2006-01-02", string(raw))
	*d = date(t)
	return err
}

var _ encoding.TextUnmarshaler = (*date)(nil)
var _ json.Unmarshaler = (*Prices)(nil)

func (p *Prices) UnmarshalJSON(raw []byte) error {
	if *p == nil {
		*p = make(map[time.Time]model.Price)
	}

	var resp struct {
		TimeSeries map[date]struct {
			Open  string `json:"1. open"`
			High  string `json:"2. high"`
			Low   string `json:"3. low"`
			Close string `json:"4. close"`
		} `json:"Time Series (Daily)"`
		ErrMsg *string `json:"Error Message"`
	}
	err := json.Unmarshal(raw, &resp)
	if err != nil {
		return err
	}

	if resp.ErrMsg != nil {
		if strings.HasPrefix(*resp.ErrMsg, "Invalid API call") {
			return model.BadRequest{"Invalid API call. E.g. ticker is wrong."}
		}
		return errors.New(*resp.ErrMsg)
	}

	for d, price := range resp.TimeSeries {
		(*p)[time.Time(d)] = model.Price{
			Open:  price.Open,
			High:  price.High,
			Low:   price.Low,
			Close: price.Close,
		}
	}
	return nil
}
