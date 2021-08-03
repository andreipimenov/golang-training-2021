package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/andreipimenov/golang-training-2021/internal/service"
)

type Handler struct {
	service Service
}

type Service interface {
	GetPrice(string, time.Time) (model.Price, error)
}

var _ Service = (*service.Service)(nil)

func New(srv Service) *Handler {
	return &Handler{
		service: srv,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	date := chi.URLParam(r, "date")

	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		writeResponse(w, model.BadRequest{"Invalid date format. The date should look like 2006-01-02"})
		return
	}

	price, err := h.service.GetPrice(ticker, d)
	if err != nil {
		writeResponse(w, err)
		return
	}

	writeResponse(w, price)
}
