package repository

import (
	"context"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/go-redis/cache/v8"
	"github.com/rs/zerolog"
	"time"
)

type RedisStore struct {
	*cache.Cache
	Logger *zerolog.Logger
}

func (r *RedisStore) Load(key string) (model.Price, bool) {
	var price model.Price
	err := r.Get(context.Background(), key, &price)
	if err != nil {
		return price, false
	}
	return price, true
}

func (r *RedisStore) Store(key string, value model.Price) {
	err := r.Set(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   time.Minute * 60,
	})
	if err != nil {
		r.Logger.Error().Err(err)
	}
}
