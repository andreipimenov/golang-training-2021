package main

type FormattedResponse struct {
	Ticker        string
	Highest_price float64
	Lowest_price  float64
	Avg_price     float64
}

// Should add custom unmarshal
//func (sr *StockRequest) UnmarshalJSON(b []byte) error {
//}
