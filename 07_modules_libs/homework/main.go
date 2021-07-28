package main

import (
  "encoding/json"
	"context"
	"fmt"
	"log"
	"net/http"
  "io/ioutil"
	"os"
	"os/signal"
  "strings"
	"syscall"
	"time"

  "github.com/go-chi/chi/v5"
)

type Stock struct {
	Ticker  string `json:"ticker"`
	Date  string `json:"date"`
  ClosePrice string `json:"close_price"`
}

type Field struct {
     ClosePrice string `json:"4. close"`

}


type AlphaVantage struct {
  Daily  map[string]Field `json:"Time Series (Daily)"`
  
}



func ping(w http.ResponseWriter, r *http.Request) {

 

  key := "R2Y0YN0KC6E92ZNN"
    ticker := chi.URLParam(r, "ticker")
  date := chi.URLParam(r, "date")

   url := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=%v&outputsize=full&apikey=%v", strings.ToUpper(ticker), key)

  fmt.Println(url)

  resp, err := http.Get(url)

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  raw := []byte(body)
	s2 := AlphaVantage{
  }

	json.Unmarshal(raw, &s2)

  

 
  s := Stock{
    Ticker: ticker,
    Date: date,
    ClosePrice: s2.Daily[date].ClosePrice,

  }

  b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(b))
  
}

func main() {


  r := chi.NewRouter()


	srv := http.Server{
		Addr: ":8080",
    Handler: r,
	}

	r.MethodFunc(http.MethodGet, "/price/{ticker}/date/{date}", ping)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	go func() {
		log.Println("Server is listening on :8080")
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-shutdown

	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped gracefully")

}
