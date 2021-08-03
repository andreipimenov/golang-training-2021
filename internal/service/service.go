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
	Load(string) (model.Ticker, bool)
	Store(string, model.Ticker)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (s *Service) GetPrice(ticker string, date time.Time) (*model.Price, error) {

	tickerInfo, ok := s.repo.Load(ticker)
	var sar stockAPIResponse
	if !ok {
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

		err = json.Unmarshal(b, &sar)
		if err != nil {
			return nil, err
		}

		tickerInfo.History = sar.History
		tickerInfo.LastRefreshed = sar.LastRefreshed
		tickerInfo.Name = sar.Name
		s.repo.Store(ticker, tickerInfo)
	}

	//The information in the API used is updated at the end of the day.
	//Here we check that the current or future day is not requested
	if tickerInfo.LastRefreshed.Unix() < date.Unix() {
		return nil, fmt.Errorf("Predictor mode is under development. The release date is never")
	}

	//The exchange may not work on some days, for example, on holidays.
	//Therefore, you need to check whether there is information on the specified date
	p, ok := tickerInfo.History[date]
	if !ok {
		return nil, fmt.Errorf("There is no information on this date")
	}

	return &p, nil
}
