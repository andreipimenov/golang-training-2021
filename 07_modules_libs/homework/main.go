package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

func logErr(err error) {
	log.Println("Error: ", err)
}

func connectionsClosedForServer(server *http.Server) chan struct{} {
	connectionsClosed := make(chan struct{})
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		defer signal.Stop(shutdown)
		<-shutdown

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		log.Println("Closing connections")
		if err := server.Shutdown(ctx); err != nil {
			logErr(err)
		}
		close(connectionsClosed)
	}()
	return connectionsClosed
}

func serveTickerStat(w http.ResponseWriter, ticker string) {
	log.Println("Serving stat of " + ticker)

	reportICE := func(err error) {
		w.WriteHeader(http.StatusInternalServerError)
		logErr(err)
	}

	res, err := http.Get("https://query1.finance.yahoo.com/v8/finance/chart/?symbol=" + ticker + "&period1=0&period2=9999999999&interval=3mo")
	if err != nil {
		reportICE(err)
		return
	}
	defer res.Body.Close()

	var yahooStat struct {
		Chart struct {
			Result []struct {
				Indicators struct {
					Quote []struct {
						Low   []float64
						High  []float64
						Open  []float64
						Close []float64
					}
				}
			}
			Error interface{}
		}
	}
	err = json.NewDecoder(res.Body).Decode(&yahooStat)
	if err != nil {
		reportICE(err)
		return
	}
	if yahooStat.Chart.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	quote := yahooStat.Chart.Result[0].Indicators.Quote[0]

	stat := struct {
		Ticker       string  `json:"ticker"`
		LowestPrice  float64 `json:"lowest_price"`
		HighestPrice float64 `json:"highest_price"`
		AvgPrice     float64 `json:"avg_price"`
	}{
		Ticker:       ticker,
		LowestPrice:  float64(math.MaxFloat64),
		HighestPrice: float64(-math.MaxFloat64),
	}
	for i := range quote.Low {
		if stat.LowestPrice > quote.Low[i] {
			stat.LowestPrice = quote.Low[i]
		}
		if stat.HighestPrice < quote.High[i] {
			stat.HighestPrice = quote.High[i]
		}
		stat.AvgPrice += (quote.Open[i] + quote.Close[i]) / 2.0 / float64(len(quote.Low))
	}

	err = json.NewEncoder(w).Encode(stat)
	if err != nil {
		reportICE(err)
		return
	}
}

func main() {
	router := chi.NewRouter()
	addr := ":8080"
	server := http.Server{Addr: addr, Handler: router}

	router.Get("/price/{ticker}/stat", func(w http.ResponseWriter, req *http.Request) {
		ticker := chi.URLParam(req, "ticker")
		serveTickerStat(w, ticker)
	})

	connectionsClosed := connectionsClosedForServer(&server)
	log.Println("Server is listening on " + addr)
	if e := server.ListenAndServe(); e != http.ErrServerClosed {
		logErr(e)
	}
	<-connectionsClosed
}
