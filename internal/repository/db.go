package repository

import (
	"database/sql"
	cfg "github.com/andreipimenov/golang-training-2021/internal/config"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type DB struct {
	*sql.DB
}

func NewDB(db *sql.DB) *DB {
	return &DB{db}
}

func (db *DB) Load(key string) (model.Price, bool) {
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

func (db *DB) Store(key string, value model.Price) {
	ticker, date := splitKey(key)
	_, err := db.Exec("INSERT INTO prices (ticker, price_date, open, high, low, close) VALUES ($1, $2, $3, $4, $5, $6)", ticker, date, value.Open, value.High, value.Low, value.Close)
	if err != nil {
		log.Err(err).Msg("Error on store")
	}
}

func splitKey(key string) (string, string) {
	x := strings.Split(key, "_")
	if len(x) != 2 {
		return "", ""
	}
	return x[0], x[1]
}

func dbInit() *sql.DB {
	db, err := sql.Open(cfg.Get().DbDriverName, cfg.Get().DBConnString)
	if err != nil {
		log.Fatal().Err(err).Msg("DB initializing error")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("DB pinging error")
	}
	return db
}

func dbMigrations(db *sql.DB) {
	log.Debug().Msg("DB migrations starting")
	_, err := db.Exec(cfg.Get().DbMigrations)
	if err != nil {
		log.Print(err)
	}
}

func closeDb(db *sql.DB) func() {
	return func() {
		log.Debug().Msg("Close Db")
		err := db.Close()
		if err != nil {
			log.Debug().Err(err)
		}
	}
}
