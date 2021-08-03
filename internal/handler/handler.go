package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type Handler struct {
	service Service
}

type Service interface {
	GetPrice(string, time.Time) (*model.Price, error)
}

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
		writeResponse(w, http.StatusBadRequest, model.Error{"Bad request"})
		return
	}

	if d.Weekday() == time.Sunday || d.Weekday() == time.Saturday {
		writeResponse(w, http.StatusBadRequest, model.Error{"A day off is specified. The exchange does not work on weekends"})
		return
	}

	price, err := h.service.GetPrice(ticker, d)
	if err != nil {
		fmt.Println(err)
		writeResponse(w, http.StatusInternalServerError, model.Error{fmt.Sprintf("%v", err)})
		return
	}

	writeResponse(w, http.StatusOK, price)
}
