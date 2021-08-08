package service_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	m "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

const (
	ticker             = "AAPL"
	validApiKey        = "123123"
	invalidTicker      = "DIMMUBORGIR"
	date               = "2021-07-26"
	weekendDate        = "2021-07-25"
	futureDate         = "2099-01-01"
	dateLayout         = "2006-01-02"
	invalidApiKey      = ""
	invalidApiKeyMsg   = "apikey is invalid or missing"
	invalidTickerMsg   = "failed retrieving data from external API: Invalid API call."
	tooManyRequestsMsg = "failed retrieving data from external API: Thank you"
	futureDateMsg      = "failed to find price of"
	emptyPriceMsg      = "failed to find price of"
)

// Valid AAPL price for 2021-07-26
var validPrice = model.Price{
	Open:  "148.2700",
	High:  "149.8300",
	Low:   "147.7000",
	Close: "148.9900",
}

type serviceTestSuite struct {
	suite.Suite
	s          *service.Service
	repoMock   *mock.Repository
	clientMock *mock.HTTPClient
}

func (suite *serviceTestSuite) SetupTest() {
	repo := &mock.Repository{}
	client := &mock.HTTPClient{}
	s := service.New(&zerolog.Logger{}, repo, validApiKey, client)
	suite.s = s
	suite.repoMock = repo
	suite.clientMock = client
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestValidPriceFromRepo() {
	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(validPrice, true)

	res, err := suite.s.GetPrice(ticker, d)

	suite.NoError(err)
	suite.Equal(res, &validPrice)
}

func (suite *serviceTestSuite) TestInvalidApiKey() {
	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	resp, err := httpRespFromFile("error_api_key.json", 200)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(model.Price{}, false)
	suite.repoMock.On("Store", key(ticker, d), m.Anything).Once().Return(true)
	suite.clientMock.On("Do", m.Anything).Once().Return(resp, nil)

	_, err = suite.s.GetPrice(ticker, d)
	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", invalidApiKeyMsg), err)
}

func (suite *serviceTestSuite) TestInvalidTicker() {
	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	resp, err := httpRespFromFile("error_invalid_api_call.json", 200)
	suite.NoError(err)

	suite.repoMock.On("Load", key(invalidTicker, d)).Once().Return(model.Price{}, false)
	suite.repoMock.On("Store", key(ticker, d), m.Anything).Once().Return(true)
	suite.clientMock.On("Do", m.Anything).Once().Return(resp, nil)

	_, err = suite.s.GetPrice(invalidTicker, d)

	suite.Error(err)
	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", invalidTickerMsg), err)
}

func (suite *serviceTestSuite) TestValidPriceFromAV() {
	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	resp, err := httpRespFromFile("full_response.json", 200)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(model.Price{}, false)
	suite.repoMock.On("Store", key(ticker, d), m.Anything).Once().Return(true)
	suite.clientMock.On("Do", m.Anything).Once().Return(resp, nil)

	res, err := suite.s.GetPrice(ticker, d)

	suite.NoError(err)
	suite.Equal(res, &validPrice)
}

func (suite *serviceTestSuite) TestFutureDate() {
	d, err := time.Parse(dateLayout, futureDate)
	suite.NoError(err)

	resp, err := httpRespFromFile("full_response.json", 200)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(model.Price{}, false)
	suite.repoMock.On("Store", key(ticker, d), m.Anything).Once().Return(true)
	suite.clientMock.On("Do", m.Anything).Once().Return(resp, nil)

	_, err = suite.s.GetPrice(ticker, d)

	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", futureDateMsg), err)
}

func (suite *serviceTestSuite) TestEmptyPriceFromAV() {
	d, err := time.Parse(dateLayout, weekendDate)
	suite.NoError(err)

	resp, err := httpRespFromFile("full_response.json", 200)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(model.Price{}, false)
	suite.repoMock.On("Store", key(ticker, d), m.Anything).Once().Return(true)
	suite.clientMock.On("Do", m.Anything).Once().Return(resp, nil)

	_, err = suite.s.GetPrice(ticker, d)
	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", emptyPriceMsg), err)
}

func (suite *serviceTestSuite) TestTooManyRequests() {
	d, err := time.Parse(dateLayout, weekendDate)
	suite.NoError(err)

	resp, err := httpRespFromFile("error_responses_limit.json", 200)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(model.Price{}, false)
	suite.repoMock.On("Store", key(ticker, d), m.Anything).Once().Return(true)
	suite.clientMock.On("Do", m.Anything).Once().Return(resp, nil)

	_, err = suite.s.GetPrice(ticker, d)
	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", tooManyRequestsMsg), err)
}

func key(ticker string, date time.Time) string {
	return fmt.Sprintf("%s_%s", ticker, date.Format("2006-01-02"))
}

// httpRespFromFile makes a *http.Response with body containing a data from
// specified file
func httpRespFromFile(filename string, respCode int) (*http.Response, error) {
	res := new(http.Response)
	res.StatusCode = respCode
	fileContent, err := os.ReadFile(fmt.Sprintf("testdata/%v", filename))
	if err != nil {
		return nil, err
	}
	res.Body = io.NopCloser(strings.NewReader(string(fileContent)))
	return res, err
}
