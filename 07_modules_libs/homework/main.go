package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	r := chi.NewRouter()
	r.Get("/price/{ticker}/date/{date}", stockHandler)
	port:=":8080"
	srv:=http.Server{
		Addr: port,
		Handler: r,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		log.Println("Server is listening on ", port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-shutdown
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped gracefully")
}

func stockHandler(w http.ResponseWriter,r *http.Request)  {

	ticker:=chi.URLParam(r,"ticker")
	date:=chi.URLParam(r,"date")
	token:="3UGTIIOY1CEXRDCE"

	urlBuilder :=fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s", ticker, token)

	response,err:=http.Get(urlBuilder)

	if err!=nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}

	defer response.Body.Close()
	dataRow,err:=ioutil.ReadAll(response.Body)
	if err!=nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}
	var dataParsed interface{}
	err = json.Unmarshal(dataRow, &dataParsed)
	if err!=nil {
		log.Println(err)
		writeResponse(w, http.StatusBadRequest, Error{err.Error()})
		return
	}
	dataDaily,ok:=dataParsed.(map[string]interface{})["Time Series (Daily)"]
	if !ok {
		writeResponse(w, http.StatusNotFound, Error{"There is no information"})
		return
	}
	dataExactDay,ok:=dataDaily.(map[string]interface{})[date]
	if !ok {
		writeResponse(w, http.StatusNotFound, Error{"There is no such date"})
		return
	}
	closeVar,ok:=dataExactDay.(map[string]interface{})["4. close"]
	if !ok {
		writeResponse(w, http.StatusNotFound, Error{"No information about price"})
		return
	}
	stock:=StockDaily{
		Close: closeVar.(string),
		Date: date,
		Ticker: ticker,
	}
	writeResponse(w,http.StatusOK,stock)
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	b, _ := json.Marshal(v)
	w.WriteHeader(code)
	w.Write(b)
}

type StockDaily struct {
	Ticker	string	`json:"ticker"`
	Date	string	`json:"date"`
	Close 	string	`json:"close_price"`
}
type Error struct {
	Error string
}