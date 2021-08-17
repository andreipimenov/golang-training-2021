package repository

import (
	"context"
	"encoding/json"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type DB struct {
	*redis.Client
	Logger *zerolog.Logger
}

func (db *DB) Load(key string) (model.Price, bool) {
	var price model.Price
	ctx := context.TODO()
	val, err := db.Get(ctx, key).Result()
	err = json.Unmarshal([]byte(val), &price)
	switch {
	case err == redis.Nil:
		return price, false
	case err != nil:
		return price, false
	case val == "":
		return price, false
	}
	return price, true
}

func (db *DB) Store(key string, value model.Price) {
	val, err := json.Marshal(value)
	if err != nil {
		db.Logger.Error().Err(err)
	}
	ctx := context.TODO()
	db.Set(ctx, key, val, 0)
}
