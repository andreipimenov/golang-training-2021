package main

import "sync"

type connectField struct {
	DailyRes map[string]Field `json:"Time Series (Daily)"`
}
type Field struct {
	Price string `json:"4. close"`
}
type UserHandler struct {
	data sync.Map
}
type respStr struct {
	Ticker         string  `json:"ticker"`
	PercentageDiff float64 `json:"diff"`
	FirstDate      string  `json:"first_date"`
	SecondDate     string  `json:"second_date"`
}
type Error struct {
	Error string
}
