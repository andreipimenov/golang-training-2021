package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

const (
	Path = "/price/{ticker}/{date}"
)

type Handler struct {
	logger  *zerolog.Logger
	service Service
}

//go:generate mockery --output $PWD/internal/mock --outpkg mock --name=Service
type Service interface {
	GetPrice(string, time.Time) (*model.Price, error)
}

func New(logger *zerolog.Logger, srv Service) *Handler {
	return &Handler{
		logger:  logger,
		service: srv,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	date := chi.URLParam(r, "date")

	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid incoming date parameter")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
		return
	}

	price, err := h.service.GetPrice(ticker, d)
	if err != nil {
		h.logger.Error().Err(err).Msg("GetPrice method error")
		writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Internal server error"})
		return
	}

	writeResponse(w, http.StatusOK, price)
}
