package service_test

import (
	"fmt"
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
	ticker           = "AAPL"
	invalidTicker    = "DIMMUBORGIR"
	strangeTicker    = "\u007f"
	date             = "2021-07-26"
	weekendDate      = "2021-07-25"
	futureDate       = "2099-01-01"
	dateLayout       = "2006-01-02"
	invalidDate      = "2021-07-269999999999"
	apiFormat        = "http://127.0.0.1:8080/price/%s/%s"
	validApiKey      = "123123"
	invalidApiKey    = ""
	invalidApiKeyMsg = "apikey is invalid or missing"
	invalidTickerMsg = "failed retrieving data from external API: Invalid API call."
	futureDateMsg    = "failed to find price of"
	strangeTickerMsg = "invalid control character in URL"
)

// Valid AAPL price for 2021-07-26
var validPrice = model.Price{
	Open:  "148.2700",
	High:  "149.8300",
	Low:   "147.7000",
	Close: "148.9900",
}

// "valid" price for 2021-07-25 (sunday)
var emptyPrice = model.Price{
	Open:  "",
	High:  "",
	Low:   "",
	Close: "",
}

type serviceTestSuite struct {
	suite.Suite
	s        *service.Service
	repoMock *mock.Repository
}

func (suite *serviceTestSuite) SetupTest() {
	repo := &mock.Repository{}
	s := service.New(&zerolog.Logger{}, repo, validApiKey)
	suite.s = s
	suite.repoMock = repo
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) TestInvalidApiKey() {
	suite.s = service.New(&zerolog.Logger{}, suite.repoMock, invalidApiKey)

	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(model.Price{}, false)

	_, err = suite.s.GetPrice(ticker, d)
	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", invalidApiKeyMsg), err)
}

func (suite *serviceTestSuite) TestEmptyPrice() {
	// dunno if it has any sense...
	d, err := time.Parse(dateLayout, weekendDate)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(emptyPrice, true)

	res, err := suite.s.GetPrice(ticker, d)
	// We'll generate an error only in http_handler in this case
	suite.NoError(err)
	suite.Equal(res, &emptyPrice)
}

func (suite *serviceTestSuite) TestInvalidTicker() {
	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	suite.repoMock.On("Load", key(invalidTicker, d)).Once().Return(model.Price{}, false)

	_, err = suite.s.GetPrice(invalidTicker, d)

	suite.Error(err)
	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", invalidTickerMsg), err)
}

func (suite *serviceTestSuite) TestValidPriceFromAV() {
	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(model.Price{}, false)
	suite.repoMock.On("Store", key(ticker, d), m.Anything).Once().Return(true)

	res, err := suite.s.GetPrice(ticker, d)

	suite.NoError(err)
	suite.Equal(res, &validPrice)
}

func (suite *serviceTestSuite) TestValidPriceFromRepo() {
	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(validPrice, true)

	res, err := suite.s.GetPrice(ticker, d)

	suite.NoError(err)
	suite.Equal(res, &validPrice)
}

func (suite *serviceTestSuite) TestFutureDate() {
	d, err := time.Parse(dateLayout, futureDate)
	suite.NoError(err)

	suite.repoMock.On("Load", key(ticker, d)).Once().Return(model.Price{}, false)

	_, err = suite.s.GetPrice(ticker, d)

	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", futureDateMsg), err)
}

func (suite *serviceTestSuite) TestUnparsableUrl() {
	d, err := time.Parse(dateLayout, date)
	suite.NoError(err)

	suite.repoMock.On("Load", key(strangeTicker, d)).Once().Return(model.Price{}, false)

	_, err = suite.s.GetPrice(strangeTicker, d)

	suite.Assert().Regexp(fmt.Sprintf(".*%v.*", strangeTickerMsg), err)
}

func key(ticker string, date time.Time) string {
	return fmt.Sprintf("%s_%s", ticker, date.Format("2006-01-02"))
}
