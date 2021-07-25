package main

// TickerDiff is a struct for our cliens requests
type TickerDiff struct {
	Ticker     string `json:"ticker"`
	FirstDate  string `json:"first_date"`
	SecondDate string `json:"second_date"`
	Format     string `json:"-"`
	Diff       string `json:"percentage_diff"`
}

// TimeSeriesDaily represents dayli data for stock
// example: https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=AAPL&outputsize=compact&apikey=demo
type TimeSeriesDaily struct {
	MetaData       map[string]string    `json:"Meta Data"`
	TimeSeriesData map[string]DailyData `json:"Time Series (Daily)"`
	ErrorMessage   string               `json:"Error Message"`
}

type DailyData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type APIError struct {
	Error string `json:"error"`
}
