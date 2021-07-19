package cache

import (
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type Item struct {
	value      		interface{}
	creationTime	time.Time
	expirationTime	int64
}

type SimpleCache struct {
	sync.Mutex
	items             map[string]Item
}

func (c *SimpleCache) Set(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	creation_time := time.Now(); c.items[key] = Item{
		value:			value,
		creationTime:	creation_time,
		expirationTime:	creation_time.Add(ttl).Unix(),
	}
}
func (c *SimpleCache) Get(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	item, ok := c.items[key]
	return item.value, ok
}

func (c *SimpleCache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func NewCache() *SimpleCache {
	cache := &SimpleCache{items: make(map[string]Item)}
	go cache.startGC()
	return cache
}

func (c *SimpleCache) startGC() {
	for ; c.items != nil; <-time.After(time.Second) {
		if keys := c.getExpiredKeys(); len(keys) != 0 {
			for _, key := range keys {
				c.Delete(key)
			}
		}
	}
}

func (c *SimpleCache) getExpiredKeys() (keys []string) {
	c.Lock()
	defer c.Unlock()
	for key, item := range c.items {
		if time.Now().Unix() >= item.expirationTime {
			keys = append(keys, key)
		}
	}
	return
}
