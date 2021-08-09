package service_test

import (
	"bytes"
	"fmt"
	"github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/repository"
	"github.com/andreipimenov/golang-training-2021/internal/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var (
	price = model.Price{
		Open:  "99.9",
		High:  "99.9",
		Low:   "99.9",
		Close: "99.9",
	}
	apiKeyValid     = "asdasd"
	apiKeyInvalid   = ""
	tickerValid     = "AAPL"
	tickerInvalid   = "ASDASD"
	dateValid       = "2021-07-26"
	dateLayoutValid = "2006-01-02"
	stockApiFormat  = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s"
)

type serviceTestSuite struct {
	suite.Suite
}

func TestService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestHitCache() {
	p, err := service.New(
		&zerolog.Logger{},
		mock.NewMockRepository().WithMethods(mock.NewMethod().
			WithName("Load").
			WithArguments(service.Key(tickerValid, dateValidParsed(suite))).
			WithReturns(price, true)).Build(),
		apiKeyValid,
		&http.Client{}).
		GetPrice(tickerValid, dateValidParsed(suite))
	suite.NoError(err)
	suite.Equal(*p, price)
}

func dateValidParsed(suite *serviceTestSuite) time.Time {
	d, err := time.Parse(dateLayoutValid, dateValid)
	suite.NoError(err)
	return d
}

func (suite *serviceTestSuite) TestInvalidTicker() {
	_, err := service.New(
		&zerolog.Logger{},
		repository.New(),
		apiKeyValid,
		mock.NewMockHTTPClient().WithMethods(mock.NewMethod().
			WithName("Do").
			WithArguments(requestTickerInvalid(suite)).
			WithReturns(nil, fmt.Errorf("err"))).Build()).
		GetPrice(tickerInvalid, dateValidParsed(suite))
	suite.Error(err)
	suite.EqualError(err, "err")
}

func requestTickerInvalid(suite *serviceTestSuite) *http.Request {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockApiFormat, tickerInvalid, apiKeyValid), nil)
	suite.NoError(err)
	return r
}

func requestApiKeyInvalid(suite *serviceTestSuite) *http.Request {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockApiFormat, tickerValid, apiKeyInvalid), nil)
	suite.NoError(err)
	return r
}

func (suite *serviceTestSuite) TestInvalidApiKey() {
	_, err := service.New(&zerolog.Logger{},
		repository.New(),
		apiKeyInvalid,
		mock.NewMockHTTPClient().WithMethods(mock.NewMethod().
			WithName("Do").
			WithArguments(requestApiKeyInvalid(suite)).
			WithReturns(nil, fmt.Errorf("err"))).Build()).
		GetPrice(tickerValid, dateValidParsed(suite))
	suite.Error(err)
	suite.EqualError(err, "err")
}

func (suite *serviceTestSuite) TestGetFromAPISuccess() {
	p, err := service.New(
		&zerolog.Logger{},
		repository.New(),
		apiKeyValid,
		mock.NewMockHTTPClient().WithMethods(mock.NewMethod().
			WithName("Do").
			WithArguments(requestValid(suite)).
			WithReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(jsonValid()))),
			}, nil)).Build()).
		GetPrice(tickerValid, dateValidParsed(suite))
	suite.NoError(err)
	suite.Equal(*p, price)
}

func jsonValid() string {
	return `{
  "Meta Data": {
    "1. Information": "Daily Prices (open, high, low, close) and Volumes",
    "2. Symbol": "AAPL",
    "3. Last Refreshed": "2021-08-11",
    "4. Output Size": "Compact",
    "5. Time Zone": "US/Eastern"
  },
  "Time Series (Daily)": {
    "2021-07-26": {
      "1. open": "99.9",
      "2. high": "99.9",
      "3. low": "99.9",
      "4. close": "99.9",
      "5. volume": "99.9"
    }
  }
}`
}

func requestValid(suite *serviceTestSuite) *http.Request {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockApiFormat, tickerValid, apiKeyValid), nil)
	suite.NoError(err)
	return r
}
