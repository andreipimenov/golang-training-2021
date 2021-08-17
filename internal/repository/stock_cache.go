package repository

import (
	"sync"

	"github.com/andreipimenov/golang-training-2021/internal/model"
)

type Cache struct {
	data sync.Map
}

func NewCache() *Cache {
	return &Cache{data: sync.Map{}}
}

func (c *Cache) Load(key string) (model.Price, bool) {
	value, ok := c.data.Load(key)
	if !ok {
		return model.Price{}, false
	}
	p, ok := value.(model.Price)
	return p, ok
}

func (c *Cache) Store(key string, value model.Price) {
	c.data.Store(key, value)
}
