package repository

import (
	"database/sql"
	"strings"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type Stock struct {
	*sql.DB
}

func NewStock(db *sql.DB) *Stock {
	return &Stock{db}
}

func (db *Stock) Load(key string) (model.Price, bool) {
	var open, high, low, close string
	ticker, date := splitKey(key)
	err := db.QueryRow("SELECT open, high, low, close FROM prices WHERE price_date = $1 AND ticker = $2", date, ticker).Scan(&open, &high, &low, &close)
	if err != nil {
		return model.Price{}, false
	}
	return model.Price{
		Open:  open,
		High:  high,
		Low:   low,
		Close: close,
	}, true
}

func (db *Stock) Store(key string, value model.Price) {
	ticker, date := splitKey(key)
	db.Exec("INSERT INTO prices (ticker, price_date, open, high, low, close) VALUES ($1, $2, $3, $4, $5, $6)", ticker, date, value.Open, value.High, value.Low, value.Close)
}

func splitKey(key string) (string, string) {
	x := strings.Split(key, "_")
	if len(x) != 2 {
		return "", ""
	}
	return x[0], x[1]
}
