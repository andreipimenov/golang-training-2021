package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	m "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
)

const (
	validTicker = "AAPL"
	dateString  = "2021-07-26"
)

var date, _ = time.Parse("2006-01-02", dateString)
var validKeyRepo = key(validTicker, date)

var validResponseJSON = fmt.Sprintf(
	`{
			"Time Series (Daily)": {
				"%s": {
					"1. open": "%s",
					"2. high": "%s",
					"3. low": "%s",
					"4. close": "%s"
				}
			}
		}`,
	dateString, price.Open, price.High, price.Low, price.Close)

var price = model.Price{
	Open:  "40.0",
	High:  "40.0",
	Low:   "40.0",
	Close: "40.0",
}

func requestChecker(ticker string) func(req *http.Request) bool {
	return func(req *http.Request) bool {
		return strings.Contains(req.URL.String(), ticker)
	}
}

type serviceTestSuite struct {
	suite.Suite
	repoMock   *mock.Repository
	clientMock *mock.HTTPClient
	service    Service
	apiKey     string
}

func (suite *serviceTestSuite) SetupTest() {
	suite.repoMock = &mock.Repository{}
	suite.clientMock = &mock.HTTPClient{}
	suite.service = *New(&zerolog.Logger{}, suite.clientMock, suite.repoMock, suite.apiKey)
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestServiceWithCache() {
	suite.repoMock.On("Load", validKeyRepo).Once().Return(price, true)
	p, err := suite.service.GetPrice(validTicker, date)
	suite.NoError(err)
	suite.Equal(price, *p)
}

func (suite *serviceTestSuite) TestServiceWithoutCache() {

	recorder := httptest.NewRecorder()
	recorder.Code = http.StatusOK
	recorder.Body.WriteString(validResponseJSON)

	suite.repoMock.On("Load", validKeyRepo).Once().Return(price, false)
	suite.clientMock.
		On("Do",
			m.MatchedBy(requestChecker(validTicker))).
		Once().Return(recorder.Result(), nil)
	suite.repoMock.On("Store", validKeyRepo, price)
	gotPrice, err := suite.service.GetPrice(validTicker, date)
	suite.NoError(err)
	suite.Equal(price, *gotPrice)
}

func (suite *serviceTestSuite) TestInvalidResponseJSON() {

	jsonFields := [4]string{
		fmt.Sprintf(`"1. open": "%s",`, price.Open),
		fmt.Sprintf(`"2. high": "%s",`, price.High),
		fmt.Sprintf(`"3. low": "%s",`, price.Low),
		fmt.Sprintf(`"4. close": "%s"`, price.Close),
	}

	replace := func(index int) string {
		if index == len(jsonFields) {
			return ","
		}
		return ""
	}

	for index, value := range jsonFields {
		resJson := strings.Replace(validResponseJSON, value, replace(index), 1)

		recorder := httptest.NewRecorder()
		recorder.Code = http.StatusOK
		recorder.Body.WriteString(resJson)

		suite.repoMock.On("Load", validKeyRepo).Once().Return(price, false)
		suite.clientMock.
			On("Do", m.MatchedBy(requestChecker(validTicker))).
			Once().Return(recorder.Result(), nil)
		suite.repoMock.On("Store", validKeyRepo, price)
		_, err := suite.service.GetPrice(validTicker, date)
		suite.Error(err)
	}
}

func (suite *serviceTestSuite) TestHTTPClient() {
	suite.repoMock.
		On("Load", validKeyRepo).
		Once().
		Return(price, false)
	e := errors.New("Http Client Do Error")
	suite.clientMock.On("Do",
		m.MatchedBy(requestChecker(validTicker))).
		Once().Return(nil, e)
	_, err := suite.service.GetPrice(validTicker, date)
	suite.Equal(e, err)
}
