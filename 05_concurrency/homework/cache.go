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

type CacheInit struct {
	mu sync.Mutex
	v  map[string]Elem
}

type Elem struct {
	val  interface{}
	date time.Time
	ttl  int64
}

func (c *CacheInit) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v[key] = Elem{
		val:  value,
		date: time.Now(),
		ttl:  time.Now().Add(ttl).UnixNano(),
	}
}

func (c *CacheInit) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.v[key]
	return val.val, ok
}

func (c *CacheInit) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.v, key)
}

func (c *CacheInit) expiredKeys() (keys []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, it := range c.v {
		if time.Now().UnixNano() >= it.ttl {
			keys = append(keys, key)
		}
	}
	return
}

func creatingCache() *CacheInit {
	cache := &CacheInit{v: make(map[string]Elem)}
	go cache.Clean()
	return cache
}

func (c *CacheInit) Clean() {
	for ; c.v != nil; <-time.After(time.Second) {
		keys := c.expiredKeys()
		if len(keys) != 0 {
			for _, key := range keys {
				c.Delete(key)
			}
		}
	}
}

func main() {
	elems := make(map[string]Elem)
	var c CacheInit = CacheInit{
		v: elems,
	}
	c.Set("Smart", "1", 2*time.Second)
	fmt.Println(c.Get("Smart"))
	c.Set("New key", "new", 3*time.Second)
	fmt.Println(c.Get("New key"))
	c.Set("Another key", 10, 5*time.Second)
	fmt.Println(c.Get("other"))
}
