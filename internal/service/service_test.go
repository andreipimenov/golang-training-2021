package service_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/andreipimenov/golang-training-2021/internal/handler"
	mocks "github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	token      string
	service    handler.Service
	repoMock   mocks.Repository
	clientMock mocks.HTTPClient
}

func TestService(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	s.repoMock = mocks.Repository{}
	s.clientMock = mocks.HTTPClient{}
	s.service = service.New(&zerolog.Logger{}, &s.repoMock, &s.clientMock, s.token)
}

func newDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func formatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func newRepoKey(ticker string, date time.Time) string {
	return ticker + "_" + formatDate(date)
}

func newResponseWithBody(json string) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.Code = http.StatusOK
	recorder.Body.WriteString(json)
	return recorder.Result()
}

func newRequestChecker(ticker string) func(req *http.Request) bool {
	return func(req *http.Request) bool {
		return strings.Contains(req.URL.String(), ticker)
	}
}

var (
	ticker       = "GOOG"
	date         = newDate(2021, 8, 5)
	price        = model.Price{Open: "2720.5700", High: "2739.0000", Low: "2712.0000", Close: "2738.8000"}
	repoKey      = newRepoKey(ticker, date)
	responseJSON = fmt.Sprintf(
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
		formatDate(date),
		price.Open, price.High, price.Low, price.Close,
	)
	newResponse = func() *http.Response {
		return newResponseWithBody(responseJSON)
	}
	checkRequest = newRequestChecker(ticker)
)

func (s *Suite) TestUncached() {
	s.repoMock.On("Load", repoKey).Once().Return(model.Price{}, false)
	s.clientMock.
		On("Do", mock.MatchedBy(checkRequest)).
		Once().Return(newResponse(), nil)
	s.repoMock.On("Store", repoKey, price)
	gotPrice, err := s.service.GetPrice(ticker, date)
	s.NoError(err)
	s.Equal(price, *gotPrice)
}

func (s *Suite) TestCached() {
	s.repoMock.On("Load", repoKey).Return(price, true)
	gotPrice, err := s.service.GetPrice(ticker, date)
	s.NoError(err)
	s.Equal(price, *gotPrice)
}

func (s *Suite) TestErrorFromHTTPClientDo() {
	s.repoMock.On("Load", repoKey).Once().Return(model.Price{}, false)
	e := errors.New("Some error")
	s.clientMock.
		On("Do", mock.MatchedBy(checkRequest)).
		Once().Return(nil, e)
	_, err := s.service.GetPrice(ticker, date)
	s.Equal(e, err)
}

func (s *Suite) TestNoDataForDate() {
	future := newDate(9999, 0, 0)
	s.repoMock.On("Load", newRepoKey(ticker, future)).Once().Return(model.Price{}, false)
	s.clientMock.
		On("Do", mock.MatchedBy(checkRequest)).
		Once().Return(newResponse(), nil)
	_, err := s.service.GetPrice(ticker, future)
	s.Error(err)
	s.Contains(err.Error(), "failed to find price")
}

func (s *Suite) TestInvalidTicker() {
	invalidTicker := "\n"
	s.repoMock.On("Load", newRepoKey(invalidTicker, date)).Once().Return(model.Price{}, false)
	_, err := s.service.GetPrice(invalidTicker, date)
	urlErr := &url.Error{}
	s.ErrorAs(err, &urlErr)
}

func (s *Suite) TestUnknownTicker() {
	s.repoMock.On("Load", repoKey).Once().Return(model.Price{}, false)
	s.clientMock.
		On("Do", mock.MatchedBy(checkRequest)).
		Once().Return(newResponseWithBody(`{"Error Message":"Invalid API call."}`), nil)
	_, err := s.service.GetPrice(ticker, date)
	s.Contains(err.Error(), "Invalid API call.")
}

type BadBody struct{}

func (b BadBody) Read(p []byte) (n int, err error) {
	return 0, errors.New("Connection lost")
}

func (s *Suite) TestErrorDuringResponseBodyReading() {
	s.repoMock.On("Load", repoKey).Once().Return(model.Price{}, false)
	response := newResponse()
	response.Body = io.NopCloser(BadBody{})
	s.clientMock.
		On("Do", mock.MatchedBy(checkRequest)).
		Once().Return(response, nil)
	_, err := s.service.GetPrice(ticker, date)
	s.Error(err)
	s.Contains(err.Error(), "Connection lost")
}

func (s *Suite) TestBadExternalAPIResponse() {
	cases := []struct {
		name, word, replacer, errFragment string
	}{
		{"Without open", "open", "", "unexpected json"},
		{"Without high", "high", "", "unexpected json"},
		{"Without low", "low", "", "unexpected json"},
		{"Without close", "close", "", "unexpected json"},
		{"Without Time Series", "Time Series", "", "unexpected json"},
		{"Value for the date is not a dictionary",
			`"` + formatDate(date), `"1999-01-01":7,"` + formatDate(date), "unexpected json"},
		{"Invalid date format", formatDate(date), "123", "parsing time"},
		{"Response is not a dictionary", "", "7", "cannot unmarshal"},
	}
	for _, c := range cases {
		s.Run(c.name, func() {
			s.repoMock.On("Load", repoKey).Once().Return(model.Price{}, false)
			bodyJSON := c.replacer
			if c.word != "" {
				bodyJSON = strings.Replace(responseJSON, c.word, c.replacer, 1)
			}
			s.clientMock.
				On("Do", mock.MatchedBy(checkRequest)).
				Once().Return(newResponseWithBody(bodyJSON), nil)
			_, err := s.service.GetPrice(ticker, date)
			s.Error(err)
			s.Contains(err.Error(), c.errFragment)
		})
	}
}
