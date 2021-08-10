package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"github.com/andreipimenov/golang-training-2021/internal/handler"
	"github.com/andreipimenov/golang-training-2021/internal/mock"
	"github.com/andreipimenov/golang-training-2021/internal/model"
)

const (
	ticker      = "AAPL"
	date        = "2021-07-26"
	invalidDate = "2021-07-269999999999"
	apiFormat   = "http://127.0.0.1:8080/price/%s/%s"
)

type handlerTestSuite struct {
	suite.Suite
	router      *chi.Mux
	serviceMock *mock.Service
}

func (suite *handlerTestSuite) SetupTest() {
	router := chi.NewRouter()
	service := &mock.Service{}
	h := handler.New(&zerolog.Logger{}, service)
	router.Get(handler.Path, h.ServeHTTP)
	suite.router = router
	suite.serviceMock = service
}

func TestHandler(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (suite *handlerTestSuite) TestHandlerBadRequest() {
	req, err := http.NewRequest("GET", fmt.Sprintf(apiFormat, ticker, invalidDate), nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.router.ServeHTTP(rr, req)

	suite.Equal(http.StatusBadRequest, rr.Result().StatusCode)
	suite.JSONEq(`{"Error":"Bad request"}`, rr.Body.String())
}

func (suite *handlerTestSuite) TestHandlerServiceFail() {

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)

	req, err := http.NewRequest("GET", fmt.Sprintf(apiFormat, ticker, date), nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.serviceMock.On("GetPrice", ticker, d).Once().Return(nil, fmt.Errorf("error"))
	suite.router.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Result().StatusCode)
	suite.JSONEq(`{"Error":"Internal server error"}`, rr.Body.String())
}

func (suite *handlerTestSuite) TestHandlerServiceOK() {

	d, err := time.Parse("2006-01-02", date)
	suite.NoError(err)

	price := &model.Price{
		Open:  "99.9",
		High:  "99.9",
		Low:   "99.9",
		Close: "99.9",
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(apiFormat, ticker, date), nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()

	suite.serviceMock.On("GetPrice", ticker, d).Once().Return(price, nil)
	suite.router.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Result().StatusCode)
	suite.JSONEq(`{"Close":"99.9", "High":"99.9", "Low":"99.9", "Open":"99.9"}`, rr.Body.String())
}
