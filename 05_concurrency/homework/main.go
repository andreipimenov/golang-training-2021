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

type Value struct {
	V   interface{}
	TTL time.Duration
	TTD time.Time
}
type MCache struct {
	mu    sync.Mutex
	Items map[string]Value
}

func (c *MCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	c.Items[key] = Value{
		V:   value,
		TTL: ttl,
		TTD: time.Now().Add(ttl),
	}
	c.mu.Unlock()

	go func() {
		for {
			<-time.After(time.Second)
			c.mu.Lock()
			if time.Now().UnixNano() >= c.Items[key].TTD.UnixNano() {
				delete(c.Items, key)
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		}
	}()
}

func (c *MCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.Items[key]
	if !ok {
		return value.V, ok
	}
	c.Items[key] = Value{
		V:   value.V,
		TTL: c.Items[key].TTL,
		TTD: c.Items[key].TTD.Add(c.Items[key].TTL),
	}
	return value.V, ok
}

func (c *MCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Items, key)
}

func main() {
	var c Cache = &MCache{
		Items: make(map[string]Value),
	}

	c.Set("key1", "nil", time.Second*4)
	c.Set("key2", "value2", time.Second*2)
	fmt.Println(c.Get("key1"))
	fmt.Println(c.Get("key2"))
	time.Sleep(5 * time.Second)
	fmt.Println(c.Get("key1"))
	fmt.Println(c.Get("key2"))

}
