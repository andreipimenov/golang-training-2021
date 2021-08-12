package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/rs/zerolog"
)

const (
	stockAPIFormat = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s"
)

type Service struct {
	logger *zerolog.Logger
	repo   Repository
	client HTTPClient
	apiKey string
}

func New(logger *zerolog.Logger, repo Repository, apiKey string) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
		client: &http.Client{
			Timeout: time.Duration(time.Minute),
		},
		apiKey: apiKey,
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
	key := Key(ticker, date)

	logger := s.logger.With().
		Str("cache_key", key).
		Logger()

	if p, ok := s.repo.Load(key); ok {
		logger.Info().Msg("Hit cache")
		return &p, nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, s.apiKey), nil)
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

	var sar stockAPIResponse
	err = json.Unmarshal(b, &sar)
	if err != nil {
		return nil, err
	}

	p, ok := sar[date]
	if !ok {
		return nil, fmt.Errorf("failed to find price of %s by %s date", ticker, date.String())
	}

	s.repo.Store(key, p)
	logger.Info().Msg("Store to cache")

	return &p, nil
}

func Key(ticker string, date time.Time) string {
	return fmt.Sprintf("%s_%s", ticker, date.Format("2006-01-02"))
}
