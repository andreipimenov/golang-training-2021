package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	srv := http.Server{
		Addr: ":8080",
	}
	go func() {
		log.Println("Server is listening on :8080")
		http.HandleFunc("/", handler)
		CheckErr(srv.ListenAndServe())
	}()
	<-shutdown
	log.Println("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer func() {
		cancel()
	}()
	CheckErr(srv.Shutdown(ctx))
	log.Println("Server stopped gracefully")

}

func handler(w http.ResponseWriter, r *http.Request) {
	switch pathSplit(r)[1] {
	case "price":
		priceHandler(w, r)
		return
	}
	http.NotFound(w, r)
}

func pathSplit(r *http.Request) []string {
	return strings.Split(r.URL.Path, "/")
}

func priceHandler(w http.ResponseWriter, r *http.Request) {
	if len(pathSplit(r)) > 4 && pathSplit(r)[3] == "date" {
		writePrice(w, r)
		return
	}
	http.NotFound(w, r)
}

func writePrice(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(toJsonByte(struct {
		Ticker      string `json:"ticker,omitempty"`
		CloserPrice string `json:"closer_price,omitempty"`
		Date        string `json:"date,omitempty"`
	}{ticker(r), getPrice(ticker(r), date(r)), date(r)}))
	CheckErr(err)
}

func date(r *http.Request) string {
	return pathSplit(r)[4]
}

func ticker(r *http.Request) string {
	return pathSplit(r)[2]
}

func getPrice(ticker, date string) string {
	resp := Resp{}
	CheckErr(json.Unmarshal(getRespBody(ticker), &resp))
	fmt.Println(resp)
	return resp.TimeSeries[date+" 20:00:00"].Close
}

func toJsonByte(i interface{}) []byte {
	marshal, err := json.Marshal(i)
	CheckErr(err)
	return marshal
}

type Resp struct {
	MetaData   MetaData         `json:"Meta Data"`
	TimeSeries map[string]Price `json:"Time Series (60min)"`
}

type MetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	Interval      string `json:"4. Interval"`
	OutputSize    string `json:"5. Output Size"`
	TimeZone      string `json:"6. Time Zone"`
}

type Price struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

func getRespBody(ticker string) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return GetBody(DoRequest(NewRequestTicker(ticker).WithContext(ctx)))
}

func GetBody(r *http.Response) []byte {
	defer func() {
		CheckErr(r.Body.Close())
	}()
	out, err := ioutil.ReadAll(r.Body)
	CheckErr(err)
	return out
}

func NewRequestTicker(ticker string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, getUrlString(ticker), nil)
	CheckErr(err)
	return request
}

func DoRequest(req *http.Request) *http.Response {
	do, err := NewHttpClient().Do(req)
	CheckErr(err)
	return do
}

func NewHttpClient() *http.Client {
	return &http.Client{}
}

func getUrlString(ticker string) string {
	return getUrl(ticker).String()
}

func getUrl(ticker string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   "www.alphavantage.co",
		Path:   "query",
		RawQuery: url.Values{
			"function": {"TIME_SERIES_INTRADAY"},
			"symbol":   {ticker},
			"interval": {"60min"},
			"apikey":   {getApiKey()},
		}.Encode(),
	}
}

func CheckErr(e error) {
	if e != nil {
		log.Println(e)
	}
}

//TODO fix
func getApiKey() string {
	return ""
}
