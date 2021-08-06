package service_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

const (
	invalidTicker  = "ASDASD"
	ticker         = "AAPL"
	date           = "2021-07-26"
	apiKey         = "asdasd"
	invalidApiKey  = ""
	invalidDate    = "2029-07-26"
	badURL         = "dsadasd%s%s"
	stockAPIFormat = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s"
)

type serviceTestSuite struct {
	suite.Suite
	clientMock *mock.HTTPClient
	repoMock   *mock.Repository
	service    *service.Service
	router     *chi.Mux
}

func (suite *serviceTestSuite) SetupTest() {
	client := &mock.HTTPClient{}
	repo := &mock.Repository{}
	s := service.New(&zerolog.Logger{}, repo, apiKey)
	suite.clientMock = client
	suite.repoMock = repo
	suite.service = s
	suite.router = chi.NewRouter()
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestServiceRequestHitCache() {
	price := model.Price{
		Open:  "99.9",
		High:  "99.9",
		Low:   "99.9",
		Close: "99.9",
	}
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(ticker, d)
	suite.repoMock.On("Load", key).Once().Return(price, true)
	p, err := suite.service.GetPrice(ticker, d)
	suite.NoError(err)
	suite.Equal(*p, price)
}

func (suite *serviceTestSuite) TestInvalidTicker() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(badURL, ticker, apiKey), nil)
	suite.NoError(err)
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(invalidTicker, d)
	mockRetVal := model.Price{}
	suite.repoMock.On("Load", key).Once().Return(mockRetVal, false)
	mockError := "failed retrieving data from external API: Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for TIME_SERIES_DAILY."
	suite.clientMock.On("Do", req).Once().Return(nil, mockError)
	_, err = suite.service.GetPrice(invalidTicker, d)
	suite.Equal(err.Error(), mockError)
}

func (suite *serviceTestSuite) TestInvalidApiKey() {
	suite.service = service.New(&zerolog.Logger{}, suite.repoMock, invalidApiKey)
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(ticker, d)
	mockRetVal := model.Price{}
	suite.repoMock.On("Load", key).Once().Return(mockRetVal, false)
	_, err = suite.service.GetPrice(ticker, d)
	invalidApiKeyError := "failed retrieving data from external API: the parameter apikey is invalid or missing. Please claim your free API key on (https://www.alphavantage.co/support/#api-key). It should take less than 20 seconds."
	suite.Equal(err.Error(), invalidApiKeyError)
}

func (suite *serviceTestSuite) TestGetFromAPISuccess() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, apiKey), nil)
	suite.NoError(err)
	resp := &http.Response{StatusCode: http.StatusInternalServerError}
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := key(ticker, d)
	price := model.Price{
		Open:  "148.2700",
		High:  "149.8300",
		Low:   "147.7000",
		Close: "148.9900",
	}
	mockRetVal := model.Price{}
	suite.repoMock.On("Load", key).Once().Return(mockRetVal, false)
	suite.repoMock.On("Store", key, price).Once()
	suite.clientMock.On("Do", req).Once().Return(resp, nil)
	p, err := suite.service.GetPrice(ticker, d)
	suite.NoError(err)
	suite.Equal(*p, price)
}

func (suite *serviceTestSuite) TestInvalidDate() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, apiKey), nil)
	suite.NoError(err)
	resp := &http.Response{StatusCode: http.StatusInternalServerError}
	d, err := time.Parse("2006-01-02", invalidDate)
	suite.NoError(err)
	key := key(ticker, d)
	mockRetVal := model.Price{}
	suite.repoMock.On("Load", key).Once().Return(mockRetVal, false)
	suite.clientMock.On("Do", req).Once().Return(resp, nil)
	_, err = suite.service.GetPrice(ticker, d)
	suite.Equal(err, fmt.Errorf("failed to find price of %s by %s date", ticker, d.String()))
}

func key(ticker string, date time.Time) string {
	return fmt.Sprintf("%s_%s", ticker, date.Format("2006-01-02"))
}
