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

type CacheItem struct {
	value interface{}
	ttl   time.Duration
}

type MemoryCache struct {
	sync.Mutex
	data map[string]CacheItem
}

func newMemoryCache() *MemoryCache {
	cache := &MemoryCache{data: make(map[string]CacheItem)}
	cache.startValidation()
	return cache
}

func (m *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	m.Lock()
	defer m.Unlock()
	m.data[key] = CacheItem{
		value: value,
		ttl:   ttl,
	}
}

func (m *MemoryCache) Get(key string) (interface{}, bool) {
	m.Lock()
	defer m.Unlock()
	item, ok := m.data[key]
	if !ok {
		return nil, false
	} else {
		//set again time before item deleted
		item.ttl = 5 * time.Second
		return item, true
	}
}

func (m *MemoryCache) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.data, key)
}
func (c *MemoryCache) startValidation() {
	go c.Validation()
}

//find and delete old valuesfrom Cache
func (c *MemoryCache) Validation() {
	for {
		<-time.After(5 * time.Second)
		keys := c.oldValues()
		if len(keys) > 0 {
			c.deleteOldValues(keys)
		}
	}
}

//delete old values from Cashe
func (c *MemoryCache) deleteOldValues(keys []string) {
	for _, k := range keys {
		c.Delete(k)
	}
}

//find old values
func (c *MemoryCache) oldValues() (oldKeys []string) {
	for key, i := range c.data {
		if time.Now().UnixNano() > int64(i.ttl.Seconds()) {
			oldKeys = append(oldKeys, key)
		}
	}
	return oldKeys
}

func main() {
	s := newMemoryCache()
	cache := s
	cache.Set("Vladimir", "go developer", 5*time.Second)
	cache.Set("Ivan", "DevOps", 5*time.Second)
	fmt.Println(cache.Get("Vladimir"))
	cache.Delete("Ivan")
	fmt.Println(cache.Get("Ivan"))
	time.Sleep(7 * time.Second)
	fmt.Println(cache.Get("Vladimir"))
}
