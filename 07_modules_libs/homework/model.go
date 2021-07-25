package main

import "net/http"

const (
	apiKey string = "8JF7GZ8X2QU1076S"
)

// client for accessing the ALPHA VANTAGE API
type aVClient struct {
	key    string
	client *http.Client
}

type StockServer struct {
	avClient   aVClient
	httpServer http.Server
}
