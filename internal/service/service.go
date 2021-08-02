package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

const (
	stockAPIFormat = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s"
	apiKey         = "G7X892PR9Q5DC69X"
)

type Service struct {
	repo   Repository
	client HTTPClient
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
		client: &http.Client{
			Timeout: time.Duration(time.Minute),
		},
	}
}

type Repository interface {
	Load(string) (model.Price, bool)
	Store(string, model.Price)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (s *Service) GetPrice(ticker string, date time.Time) (*model.Price, error) {
	if p, ok := s.repo.Load(key(ticker, date)); ok {
		return &p, nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, apiKey), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var sar *StockAPIResponse = &StockAPIResponse{
		TimeSeriesDaily:   make(map[time.Time]model.Price),
		LastRefreshedTime: time.Time{},
	}

	err = json.Unmarshal(b, &sar)
	if err != nil {
		return nil, err
	}

	if sar.LastRefreshedTime.Before(date) {
		return nil, errUnexpectedDate
	}

	if date.Weekday().String() == "Saturday" || date.Weekday().String() == "Sunday" {
		return nil, errWeekendDate
	}

	p, ok := sar.TimeSeriesDaily[date]
	if !ok {
		return nil, errDateDoesNotExist
	}

	s.repo.Store(key(ticker, date), p)

	return &p, nil
}

func key(ticker string, date time.Time) string {
	return fmt.Sprintf("%s_%s", ticker, date.Format("2006-01-02"))
}
