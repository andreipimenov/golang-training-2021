package handler

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/models"
	"github.com/andreipimenov/golang-training-2021/internal/restapi/operations/stock"
)

type Stock struct {
	logger  *zerolog.Logger
	service StockService
}

type StockService interface {
	GetPrice(string, time.Time) (*model.Price, error)
}

func NewStock(logger *zerolog.Logger, srv StockService) *Stock {
	return &Stock{
		logger:  logger,
		service: srv,
	}
}

func (h *Stock) Handle(req stock.GetPriceParams) middleware.Responder {

	d, err := time.Parse("2006-01-02", req.Date.String())
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming date parameter")
		return stock.NewGetPriceBadRequest().WithPayload(&models.Error{"Bad request"})
	}

	price, err := h.service.GetPrice(req.Ticker, d)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetPrice method error")
		return stock.NewGetPriceInternalServerError().WithPayload(&models.Error{"Internal server error"})
	}

	return stock.NewGetPriceOK().WithPayload(&models.Price{High: price.High, Low: price.Low, Open: price.Open, Close: price.Close})
}
