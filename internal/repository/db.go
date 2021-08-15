package repository

import (
	"database/sql"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/rs/zerolog"
	"os"
)

type DB struct {
	*sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{db}
}

func (db *DB) Load(key string) (model.Price, bool) {
	var open, high, low, close string
	ticker, date := SplitKey(key)
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

func (db *DB) Store(key string, value model.Price) {
	ticker, date := SplitKey(key)
	_, err := db.Exec("INSERT INTO prices (ticker, price_date, open, high, low, close) VALUES ($1, $2, $3, $4, $5, $6)", ticker, date, value.Open, value.High, value.Low, value.Close)
	if err != nil {
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		logger.Error().Err(err).Msg("store error")
	}
}
