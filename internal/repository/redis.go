package repository

import (
	"context"
	"github.com/andreipimenov/golang-training-2021/internal/model"
	"github.com/go-redis/cache/v8"
	"time"
)

var ctx = context.Background()

type RedisDB struct {
	*cache.Cache
}

/*func NewRedisDB(db *RedisDB) *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &DB{rdb}
}*/

func (rdb *RedisDB) Load(key string) (model.Price, bool) {
	var price model.Price
	if err := rdb.Get(ctx, key, &price); err == nil {
		return price, true
	} else {
		return price, false
	}
}

func (rdb *RedisDB) Store(key string, value model.Price) {
	ctx := context.TODO()
	obj := &value

	if err := rdb.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		panic(err)
	}
}
