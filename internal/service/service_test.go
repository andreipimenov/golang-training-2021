package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

const (
	ticker          = "AAPL"
	invalidTicker   = "aar"
	date            = "2021-07-26"
	dateWithoutData = "2021-08-07"
)

type serviceTestSuite struct {
	suite.Suite
	service  *service.Service
	repoMock *mock.Repository
}

func (suite *serviceTestSuite) SetupTest() {
	repo := &mock.Repository{}
	s := service.New(&zerolog.Logger{}, repo, "1")
	suite.service = s
	suite.repoMock = repo
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestServiceRepo() {
	price := model.Price{
		Open:  "99.9",
		High:  "99.9",
		Low:   "99.9",
		Close: "99.9",
	}
	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := service.Key(ticker, d)
	suite.repoMock.On("Load", key).Once().Return(price, true)
	p, err := suite.service.GetPrice(ticker, d)
	suite.NoError(err)
	suite.Equal(*p, price)
}

func (suite *serviceTestSuite) TestServiceInvalidTicker() {

	errToReturn := "failed retrieving data from external API: Invalid API call. Please retry or visit the documentation (https://www.alphavantage.co/documentation/) for TIME_SERIES_DAILY."

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)
	key := service.Key(invalidTicker, d)
	mockRetVal := model.Price{}
	suite.repoMock.On("Load", key).Once().Return(mockRetVal, false)

	_, err = suite.service.GetPrice(invalidTicker, d)

	suite.Equal(err.Error(), errToReturn)
}

func (suite *serviceTestSuite) TestServiceDateWithoutData() {

	d, err := time.Parse("2006-01-02", dateWithoutData)
	errToReturn := fmt.Sprintf("failed to find price of %s by %s date", ticker, d)
	suite.NoError(err)
	key := service.Key(ticker, d)
	mockRetVal := model.Price{}
	suite.repoMock.On("Load", key).Once().Return(mockRetVal, false)

	_, err = suite.service.GetPrice(ticker, d)

	suite.Equal(err.Error(), errToReturn)
}

func (suite *serviceTestSuite) TestServiceAllValid() {

	d, err := time.Parse("2006-01-02", date)
	var validPrice = model.Price{
		Open:  "148.2700",
		High:  "149.8300",
		Low:   "147.7000",
		Close: "148.9900",
	}
	suite.NoError(err)
	key := service.Key(ticker, d)
	mockRetVal := model.Price{}
	suite.repoMock.On("Load", key).Once().Return(mockRetVal, false)
	suite.repoMock.On("Store", key, validPrice).Once()

	result, err := suite.service.GetPrice(ticker, d)
	suite.NoError(err)
	suite.Equal(result, &validPrice)
}
