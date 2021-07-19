package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type Item struct {
	Object         interface{}
	timeOfCreation time.Time
	TTL            time.Duration
}

func (item Item) IsExpired() bool {
	return item.timeOfCreation.Add(item.TTL).Unix() <= time.Now().Unix()
}

type InMemoryCache struct {
	sync.Mutex
	cap   uint64
	len   uint64
	items map[string]Item
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {

	c.Lock()
	c.len++
	c.items[key] = Item{
		Object:         value,
		timeOfCreation: time.Now(),
		TTL:            ttl,
	}
	if c.len > c.cap {
		for k, v := range c.items {
			if v.IsExpired() {
				delete(c.items, k)
			}
		}
	}

	c.Unlock()
}

func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.Lock()
	value, ok := c.items[key]
	if !ok {
		c.Unlock()
		return nil, false
	}
	if value.IsExpired() {
		delete(c.items, key)
		c.Unlock()
		return nil, false
	}
	c.Unlock()
	return value.Object, true
}

func (c *InMemoryCache) Delete(key string) {
	c.Lock()
	delete(c.items, key)
	c.Unlock()
}

func main() {
	c := make(map[string]Item)
	cache := InMemoryCache{
		items: c,
		cap:   3,
	}

	cache.Set("key", "Some Value", 4*time.Second)
	cache.Set("2", "Some Value2", 4*time.Second)
	cache.Set("key3", "Some Value3", 2*time.Second)
	cache.Set("key4", "Some Value4", 4*time.Second)
	time.Sleep(1 * time.Second)
	fmt.Println(cache.Get("key3"))

	for i := 0; i < 8128; i++ {
		go func(i int) {
			cache.Set("key4", "Some Value"+strconv.Itoa(i), 4*time.Second)
		}(i)

	}

	fmt.Println(cache.Get("key4"))
}