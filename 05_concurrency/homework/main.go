package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type Item struct {
	value            interface{}
	ttl              time.Duration
	invalidationTime int64
}

type CacheImplementation struct {
	sync.Mutex
	items map[string]Item
}

func (cache *CacheImplementation) Set(key string, value interface{}, ttl time.Duration) {
	cache.Lock()
	defer cache.Unlock()
	cache.items[key] = Item{
		value:            value,
		ttl:              ttl,
		invalidationTime: time.Now().Add(ttl).Unix(),
	}
}

func (cache *CacheImplementation) Get(key string) (interface{}, bool) {
	cache.Lock()
	defer cache.Unlock()
	item, ok := cache.items[key]
	if time.Now().Unix() <= item.invalidationTime {
		cache.items[key] = Item{
			value:            item.value,
			ttl:              item.ttl,
			invalidationTime: time.Now().Add(item.ttl).Unix(),
		}
		return item.value, ok
	}
	return nil, ok
}

func (cache *CacheImplementation) Delete(key string) {
	cache.Lock()
	defer cache.Unlock()
	delete(cache.items, key)
}

func (cache *CacheImplementation) DeleteInvalidValues() {
	cache.Lock()
	defer cache.Unlock()
	for key, item := range cache.items {
		if time.Now().Unix() > item.invalidationTime {
			delete(cache.items, key)
		}
	}
}

func main() {
	var cache CacheImplementation = CacheImplementation{items: make(map[string]Item)}
	cache.Set("St.Petersburg", "Russia", 3*time.Second)
	cache.Set("Gnansk", "Poland", 2*time.Second)
	cache.Set("Hamburg", "Germany", 3*time.Second)
	fmt.Println(cache.Get("St.Petersburg"))
	cache.Delete("Humburg")
}
