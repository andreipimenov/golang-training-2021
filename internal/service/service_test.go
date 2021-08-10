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
	invalidTicker  = "ABCDEFG"
	validTicker    = "AAPL"
	date           = "2021-07-26"
	apiKey         = "1234"
	invalidApiKey  = ""
	invalidDate    = "2089-12-21"
	badURL         = "abcdefg"
	stockAPIFormat = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s"
)

type serviceTestSuite struct {
	suite.Suite
	clientMock *mock.HTTPClient
	repoMock   *mock.Repository
	service    *service.Service
}

func (suite *serviceTestSuite) SetupTest() {
	client := &mock.HTTPClient{}
	repo := &mock.Repository{}
	s := service.New(&zerolog.Logger{}, repo, apiKey)
	suite.clientMock = client
	suite.repoMock = repo
	suite.service = s
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

var price = model.Price{
	Open:  "30.0",
	High:  "40.0",
	Low:   "20.0",
	Close: "35.0",
}

func (suite *serviceTestSuite) TestHitCache() {

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(validTicker, d)
	suite.repoMock.On("Load", key).Once().Return(price, true)
	p, err := suite.service.GetPrice(validTicker, d)
	suite.NoError(err)
	suite.Equal(price, *p)
}

func (suite *serviceTestSuite) TestInvalidDate() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, validTicker, apiKey), nil)
	suite.NoError(err)
	date, err := time.Parse("2006-01-02", invalidDate)
	suite.NoError(err)
	key := key(validTicker, date)
	response := &http.Response{StatusCode: http.StatusInternalServerError}
	suite.repoMock.On("Load", key).Once().Return(model.Price{}, false)
	suite.clientMock.On("Do", req).Once().Return(response, nil)
	suite.repoMock.On("Store", key, price).Return()
	resp, err := suite.service.GetPrice(validTicker, date)
	var m *model.Price
	suite.Equal(m, resp)
	suite.Equal(fmt.Errorf("failed to find price of %s by %s date", validTicker, date.String()), err)
}

func (suite *serviceTestSuite) TestInvalidTicker() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, invalidTicker, apiKey), nil)
	suite.NoError(err)
	date, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(invalidTicker, date)
	reqError := "failed retrieving data from external API: Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for TIME_SERIES_DAILY."
	suite.repoMock.On("Load", key).Once().Return(model.Price{}, false)
	suite.clientMock.On("Do", req).Once().Return(nil, reqError)
	p, err := suite.service.GetPrice(invalidTicker, date)
	var m *model.Price
	suite.Equal(reqError, err.Error())
	suite.Equal(m, p)
}

func (suite *serviceTestSuite) TestValidRequest() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, validTicker, apiKey), nil)
	suite.NoError(err)
	date, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	var validPrice = model.Price{
		Open:  "148.2700",
		High:  "149.8300",
		Low:   "147.7000",
		Close: "148.9900",
	}
	key := key(validTicker, date)
	response := &http.Response{StatusCode: http.StatusOK}
	suite.repoMock.On("Load", key).Once().Return(model.Price{}, false)
	suite.clientMock.On("Do", req).Once().Return(response, nil)
	suite.repoMock.On("Store", key, validPrice).Return()
	p, err := suite.service.GetPrice(validTicker, date)
	suite.NoError(err)
	suite.Equal(validPrice, *p)
}

func key(ticker string, date time.Time) string {
	return fmt.Sprintf("%s_%s", ticker, date.Format("2006-01-02"))
}
