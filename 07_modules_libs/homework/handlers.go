package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func DiffHandler(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	firstDate := chi.URLParam(r, "first_date")
	secondDate := chi.URLParam(r, "second_date")
	format := "compact"
	err := validateDiffDates(firstDate, secondDate)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, APIError{err.Error()})
		return
	}
	// If date(s) is before now-100d
	// if true we must get all data for ticker ('full' format)
	f, _ := time.Parse(dateLayoutISO, firstDate)
	if f.Before(time.Now().Add(-100 * 24 * time.Hour)) {
		format = "full"
	}
	td := &TickerDiff{
		Ticker:     ticker,
		FirstDate:  firstDate,
		SecondDate: secondDate,
		Format:     format,
	}
	// version with context
	ctx := r.Context()
	errChan := make(chan (error), 1)
	go GetDiffAsync(errChan, td)
	select {
	case err = <-errChan:
		if err != nil {
			writeResponse(w, http.StatusBadRequest, APIError{err.Error()})
			return
		}
		writeResponse(w, http.StatusOK, td)
		return
	case <-ctx.Done():
		log.Printf("Interrupted with err: %v", ctx.Err())
		return
	}

	// without context
	// err = GetDiff(td)
	// if err != nil {
	// 	writeResponse(w, http.StatusBadRequest, APIError{err.Error()})
	// 	return
	// }
	// writeResponse(w, http.StatusOK, td)
}
