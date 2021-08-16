package repository

import (
	"context"
	"encoding/json"

	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type RedisCache struct {
	*redis.Client
	logger *zerolog.Logger
}

func NewRedisCache(rdb *redis.Client, logger *zerolog.Logger) *RedisCache {
	return &RedisCache{rdb, logger}
}

func (c *RedisCache) Load(key string) (model.Price, bool) {
	ctx := context.Background()
	res := model.Price{}
	val, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return res, false
	}
	if err != nil {
		c.logger.Error().Err(err).Msg("error during the getting price from cache")
		return res, false
	}
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		c.logger.Error().Err(err).Msg("cannot unmarsal redis data")
	}
	return res, true
}

func (c *RedisCache) Store(key string, value model.Price) {
	ctx := context.Background()
	jsonVal, err := json.Marshal(value)
	if err != nil {
		c.logger.Error().Err(err).Msgf("cannot marshal %T object", value)
	}
	err = c.Set(ctx, key, jsonVal, 0).Err()
	if err != nil {
		c.logger.Error().Err(err).Msg("cannot cache data into redis")
	}
}
