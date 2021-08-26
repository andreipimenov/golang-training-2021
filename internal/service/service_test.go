package service_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

const (
	ticker         = "AAPL"
	invalidTicker  = "AAPLl"
	date           = "2021-08-16"
	dayOff         = "2021-08-15"
	apiKey         = "123"
	invalidApiKey  = ""
	stockAPIFormat = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s"
)

var (
	price = model.Price{
		Open:  "148.5350",
		High:  "151.1900",
		Low:   "146.4700",
		Close: "151.1200",
	}
)

type serviceTestSuite struct {
	suite.Suite
	service    *service.Service
	repoMock   *mock.Repository
	clientMock *mock.HTTPClient
}

func (suite *serviceTestSuite) SetupTest() {
	repo := &mock.Repository{}
	client := &mock.HTTPClient{}
	s := service.New(&zerolog.Logger{}, repo, client, apiKey)
	suite.service = s
	suite.repoMock = repo
	suite.clientMock = client
}

func TestService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestServiceDataFromCache() {

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)

	suite.repoMock.On("Load", service.Key(ticker, d)).Once().Return(price, true)
	p, err := suite.service.GetPrice(ticker, d)
	suite.NotNil(p)
	suite.NoError(err)
	suite.Equal(price, *p)
}

func (suite *serviceTestSuite) TestServiceDataFromAPI() {

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)

	suite.repoMock.On("Load", service.Key(ticker, d)).Once().Return(model.Price{}, false)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, apiKey), nil)
	suite.NoError(err)

	respBody := `{
		"Time Series (Daily)": {
			"2021-08-16": {
				"1. open": "148.5350",
				"2. high": "151.1900",
				"3. low": "146.4700",
				"4. close": "151.1200",
				"5. volume": "103558782"
			}
	}}`
	resp := &http.Response{}
	resp.StatusCode = http.StatusOK
	resp.Body = io.NopCloser(bytes.NewBufferString(respBody))

	suite.clientMock.On("Do", req).Once().Return(resp, nil)

	suite.repoMock.On("Store", service.Key(ticker, d), price)
	p, err := suite.service.GetPrice(ticker, d)
	suite.NotNil(p)
	suite.NoError(err)
	suite.Equal(price, *p)
}

func (suite *serviceTestSuite) TestServiceInvalidTicker() {

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)

	suite.repoMock.On("Load", service.Key(invalidTicker, d)).Once().Return(model.Price{}, false)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, invalidTicker, apiKey), nil)
	suite.NoError(err)

	invalidTickerErr := "failed retrieving data from external API: Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for TIME_SERIES_DAILY."
	respBody := `{
		"Error Message": "Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for TIME_SERIES_DAILY."
	}`
	resp := &http.Response{}
	resp.StatusCode = http.StatusOK
	resp.Body = io.NopCloser(bytes.NewBufferString(respBody))

	suite.clientMock.On("Do", req).Once().Return(resp, nil)

	p, err := suite.service.GetPrice(invalidTicker, d)
	suite.Nil(p)
	suite.Error(err)
	suite.Equal(fmt.Errorf(invalidTickerErr), err)
}

func (suite *serviceTestSuite) TestServiceDayOffDate() {

	dayOff, err := time.Parse("2006-01-02", dayOff)
	suite.NoError(err)

	suite.repoMock.On("Load", service.Key(ticker, dayOff)).Once().Return(model.Price{}, false)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, apiKey), nil)
	suite.NoError(err)

	respBody := `{
		"Time Series (Daily)": {
			"2021-08-16": {
				"1. open": "148.5350",
				"2. high": "151.1900",
				"3. low": "146.4700",
				"4. close": "151.1200",
				"5. volume": "103558782"
			}
	}}`
	resp := &http.Response{}
	resp.StatusCode = http.StatusOK
	resp.Body = io.NopCloser(bytes.NewBufferString(respBody))

	dayOffErr := `failed to find price of AAPL by 2021-08-15 00:00:00 +0000 UTC date`
	suite.clientMock.On("Do", req).Once().Return(resp, nil)

	p, err := suite.service.GetPrice(ticker, dayOff)
	suite.Nil(p)
	suite.Error(err)
	suite.Equal(fmt.Errorf(dayOffErr), err)
}

func (suite *serviceTestSuite) TestServiceInvalidStockAPIFormat() {

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)

	suite.repoMock.On("Load", service.Key(ticker, d)).Once().Return(model.Price{}, false)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, apiKey), nil)
	suite.NoError(err)
	respErr := fmt.Errorf(`Get "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=AAPL&apikey=123": dial tcp: lookup www.alphavantage.co on 172.25.96.1:53: no such host`)
	suite.clientMock.On("Do", req).Once().Return(nil, respErr)

	p, err := suite.service.GetPrice(ticker, d)
	suite.Nil(p)
	suite.Error(err)
	suite.Equal(respErr, err)
}

func (suite *serviceTestSuite) TestServiceInvalidApiKey() {

	suite.service = service.New(&zerolog.Logger{}, suite.repoMock, suite.clientMock, invalidApiKey)

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)

	suite.repoMock.On("Load", service.Key(ticker, d)).Once().Return(model.Price{}, false)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, invalidApiKey), nil)
	suite.NoError(err)

	invalidApiKeyErr := "failed retrieving data from external API: the parameter apikey is invalid or missing. Please claim your free API key on (https://www.alphavantage.co/support/#api-key). It should take less than 20 seconds."
	respBody := `{
		"Error Message": "the parameter apikey is invalid or missing. Please claim your free API key on (https://www.alphavantage.co/support/#api-key). It should take less than 20 seconds."
	}`

	resp := &http.Response{}
	resp.StatusCode = http.StatusOK
	resp.Body = io.NopCloser(bytes.NewBufferString(respBody))

	suite.clientMock.On("Do", req).Once().Return(resp, nil)

	p, err := suite.service.GetPrice(ticker, d)
	suite.Nil(p)
	suite.Error(err)
	suite.Equal(fmt.Errorf(invalidApiKeyErr), err)
}
