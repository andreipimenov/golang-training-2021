package service_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

const (
	validTicker    = "AAPL"
	validApiKey    = "123"
	validDate      = "2021-07-26"
	invalidTicker  = "AAAAAAAAAPL"
	invalidApiKey  = ""
	futureDate     = "2022-01-01"
	weekendDate    = "2021-07-25"
	stockAPIFormat = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s"
)

var (
	validPrice = model.Price{
		Open:  "148.2700",
		High:  "149.8300",
		Low:   "147.7000",
		Close: "148.9900",
	}
)

func key(ticker string, date time.Time) string {
	return fmt.Sprintf("%s_%s", ticker, date.Format("2006-01-02"))
}

type serviceTestSuite struct {
	suite.Suite
	s              *service.Service
	reposytoryMock *mock.Repository
	clientMock     *mock.HTTPClient
}

func (suite *serviceTestSuite) SetupTest() {
	reposytory := &mock.Repository{}
	client := &mock.HTTPClient{}
	s := service.New(&zerolog.Logger{}, reposytory, validApiKey)
	suite.s = s
	suite.reposytoryMock = reposytory
	suite.clientMock = client
}

func TestService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestHitCache() {
	d, err := time.Parse("2006-01-02", validDate)
	suite.NoError(err)

	key := key(validTicker, d)
	suite.reposytoryMock.On("Load", key).Once().Return(validPrice, true)
	p, err := suite.s.GetPrice(validTicker, d)
	suite.NoError(err)

	suite.Equal(*p, validPrice)
}

func (suite *serviceTestSuite) TestValidRequest() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, validTicker, validApiKey), nil)
	suite.NoError(err)

	resp := &http.Response{StatusCode: http.StatusInternalServerError}
	d, err := time.Parse("2006-01-02", validDate)
	suite.NoError(err)

	key := key(validTicker, d)

	suite.reposytoryMock.On("Load", key).Once().Return(model.Price{}, false)
	suite.reposytoryMock.On("Store", key, validPrice).Once()
	suite.clientMock.On("Do", req).Once().Return(resp, nil)
	p, err := suite.s.GetPrice(validTicker, d)
	suite.NoError(err)

	suite.Equal(*p, validPrice)
}

func (suite *serviceTestSuite) TestInvalidDate() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, validTicker, validApiKey), nil)
	suite.NoError(err)

	d, err := time.Parse("2006-01-02", validDate)
	suite.NoError(err)

	key := key(validTicker, d)
	suite.reposytoryMock.On("Load", key).Once().Return(model.Price{}, false)
	suite.reposytoryMock.On("Store", key, validPrice).Once()
	resp := &http.Response{StatusCode: http.StatusInternalServerError}
	suite.clientMock.On("Do", req).Once().Return(resp, nil)
	p, err := suite.s.GetPrice(validTicker, d)
	suite.NoError(err)

	suite.Equal(*p, validPrice)
}

func (suite *serviceTestSuite) TestWeekendDate() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, validTicker, validApiKey), nil)
	suite.NoError(err)

	d, err := time.Parse("2006-01-02", weekendDate)
	suite.NoError(err)

	key := key(validTicker, d)
	suite.reposytoryMock.On("Load", key).Once().Return(model.Price{}, false)
	suite.reposytoryMock.On("Store", key, validPrice).Once()
	resp := &http.Response{StatusCode: http.StatusInternalServerError}
	suite.clientMock.On("Do", req).Once().Return(resp, nil)
	_, err = suite.s.GetPrice(validTicker, d)

	suite.Equal(err, fmt.Errorf("failed to find price of %s by %s date", validTicker, d.String()))
}

func (suite *serviceTestSuite) TestInvalidApiKey() {
	suite.s = service.New(&zerolog.Logger{}, suite.reposytoryMock, invalidApiKey)
	d, err := time.Parse("2006-01-02", validDate)
	suite.NoError(err)

	key := key(validTicker, d)
	suite.reposytoryMock.On("Load", key).Once().Return(model.Price{}, false)
	_, err = suite.s.GetPrice(validTicker, d)
	invalidApiKeyError := "failed retrieving data from external API: the parameter apikey is invalid or missing. Please claim your free API key on (https://www.alphavantage.co/support/#api-key). It should take less than 20 seconds."

	suite.Equal(err.Error(), invalidApiKeyError)
}
