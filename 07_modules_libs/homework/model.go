package main

type Stock struct {
	Ticker       string  `json:"ticker"`
	HighestPrice float64 `json:"highest_price"`
	LowestPrice  float64 `json:"lowest_price"`
	AvgPrice     float64 `json:"avg_price"`
}

func NewStock(ticker string, avg, hi, lo float64) *Stock {
	s := new(Stock)
	s.Ticker = ticker
	s.AvgPrice = avg
	s.HighestPrice = hi
	s.LowestPrice = lo
	return s
}

type RequestJSONIntraDay struct {
	M    Meta                  `json:"Meta Data"`
	TS   map[string]SeriesData `json:"Time Series (60min)"`
	Note string                `json:"Note"`
}

type RequestJSONIndicator struct {
	M    MetaIndicator             `json:"Meta Data"`
	Data map[string]indicatorValue `json:"Technical Analysis: SMA"`
	Note string                    `json:"Note"`
}

type SeriesData struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type Meta struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	Interval      string `json:"4. Interval"`
	OutputSize    string `json:"5. Output Size"`
	TimeZone      string `json:"6. Time Zone"`
}

type MetaIndicator struct {
	Symbol        string `json:"1: Symbol"`
	Indicator     string `json:"2: Indicator"`
	LastRefreshed string `json:"3: Last Refreshed"`
	Interval      string `json:"4: Interval"`
	TimePeriod    int    `json:"5: Time Period"`
	SeriesType    string `json:"6: Series Type"`
	TimeZone      string `json:"7: Time Zone"`
}

type indicatorValue struct {
	SMA string `json:"SMA"`
}

type Error struct {
	Error string
}
