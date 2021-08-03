package model

import "time"

type Error struct {
	Error string
}

type Price struct {
	Open  string
	High  string
	Low   string
	Close string
}

//Now, when a request is made for a given ticker,
//all information about it will be stored, and not for the date
type Ticker struct {
	Name          string
	LastRefreshed time.Time
	History       map[time.Time]Price
}
