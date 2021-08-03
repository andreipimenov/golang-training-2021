package alphavantage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

func stockAPI(symbol string) string {
	return fmt.Sprintf(
		"https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s",
		url.QueryEscape(symbol),
		"G7X892PR9Q5DC69X",
	)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Repository struct {
	client HTTPClient
}

var _ service.SlowRepository = (*Repository)(nil)

func New() *Repository {
	return &Repository{client: &http.Client{
		Timeout: time.Minute,
	}}
}

func (r *Repository) GetPrice(ticker string, date time.Time) (model.Price, error) {
	req, err := http.NewRequest(http.MethodGet, stockAPI(ticker), nil)
	if err != nil {
		return model.Price{}, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return model.Price{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.Price{}, err
	}

	var prices Prices
	err = json.Unmarshal(body, &prices)
	if err != nil {
		return model.Price{}, err
	}

	p, ok := prices[date]
	if !ok {
		return model.Price{}, model.BadRequest{"No data for this date."}
	}

	return p, nil
}
