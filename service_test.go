package service

import (
	"fmt"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/handler"
	"github.com/andreipimenov/golang-training-2021/internal/mock"
)

const (
	ticker = "TSLA"
	date   = "2021-07-26"
	apiKey = "GHDKG9235GJDNDG83"
	invalidApiKey = ""
	invalidTicker = "smth"
	invalidDate = "1000-01-01"
	ApiError = "failed retrieving data from external API: the parameter apikey is invalid or missing. Please claim your free API key on (https://www.alphavantage.co/support/#api-key). It should take less than 20 seconds."
	tickerError = "failed retrieving data from external API: Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for TIME_SERIES_DAILY."
)

type serviceTestSuite struct {
	suite.Suite
	service handler.Service
	repo    mock.Repository
	client  mock.HTTPClient
	token   string
}

func TestService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) SetupTest() {
	suite.repo = mock.Repository{}
	suite.client = mock.HTTPClient{}
	suite.service = New(&zerolog.Logger{}, &suite.repo, apiKey)
}

func (suite *serviceTestSuite) TestLoadFromCache() {
	price := model.Price{
		Open:  "100.5",
		High:  "120",
		Low:   "80.5",
		Close: "105",
	}
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(ticker, d)
	suite.repo.On("Load", key).Once().Return(price, true)
	pr, err := suite.service.GetPrice(ticker, d)
	suite.NoError(err)
	suite.Equal(*pr, price)
}

func (suite *serviceTestSuite) TestInvalidApiToken()  {
	suite.service = New(&zerolog.Logger{}, &suite.repo, invalidApiKey )
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(ticker, d)
	suite.repo.On("Load", key).Once().Return(model.Price{}, false)
	_, err = suite.service.GetPrice(ticker, d)
	suite.EqualError(err, ApiError )
}

func (suite *serviceTestSuite) TestInvalidTicker() {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, invalidTicker, apiKey), nil)
	suite.NoError(err)
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(invalidTicker, d)
	suite.repo.On("Load", key).Once().Return(model.Price{}, false)
	suite.client.On("Do", request).Once().Return(nil, tickerError)
	_, err = suite.service.GetPrice(invalidTicker, d)
	suite.EqualError(err, tickerError)
}

func (suite *serviceTestSuite) TestInvalidDate() {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, apiKey), nil)
	suite.NoError(err)
	response := &http.Response{StatusCode: http.StatusInternalServerError}
	d, err := time.Parse("2006-01-02", invalidDate)
	suite.NoError(err)
	key := key(ticker, d)
	suite.repo.On("Load", key).Once().Return(model.Price{}, false)
	suite.client.On("Do", request).Once().Return(response, nil)
	_, err = suite.service.GetPrice(ticker, d)
	suite.EqualError(err, fmt.Sprintf("failed to find price of %s by %s date", ticker, d.String()))
}

func (suite *serviceTestSuite) TestSuccess() {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, apiKey), nil)
	suite.NoError(err)
	response := &http.Response{StatusCode: http.StatusInternalServerError}
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(ticker, d)
	price := model.Price{
		Open:  "650.9700",
		High:  "668.1999",
		Low:   "647.1100",
		Close: "657.6200",
	}
	suite.repo.On("Load", key).Once().Return(model.Price{}, false)
	suite.repo.On("Store", key, price).Once()
	suite.client.On("Do", request).Once().Return(response, nil)
	responsePrice, err := suite.service.GetPrice(ticker, d)
	suite.NoError(err)
	suite.Equal(*responsePrice, price)
}


