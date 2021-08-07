package main

import (
	"fmt"
	"net/http"
)

const (
	ticker         = "\u007f"
	stockAPIFormat = "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&outputsize=full&symbol=%s&apikey=%s"
)

func main() {
	_, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, "123"), nil)
	fmt.Println(err)
}
