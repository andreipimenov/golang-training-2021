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

type CacheElement struct {
	Value       interface{}
	ErasureTime time.Time
	TTL         time.Duration
}

type CacheMemory struct {
	sync.Mutex
	Memory map[string]CacheElement
}

func (cache *CacheMemory) Delete(key string) {
	cache.Lock()
	delete(cache.Memory, key)
	cache.Unlock()
}

func (cache *CacheMemory) Get(key string) (interface{}, bool) {
	cache.Lock()
	value, ok := cache.Memory[key]
	if !ok {
		cache.Unlock()
		return nil, false
	}
	//updating the time of erase
	cache.Memory[key] = CacheElement{
		Value:       value.Value,
		ErasureTime: time.Now().Add(cache.Memory[key].TTL),
		TTL:         cache.Memory[key].TTL,
	}
	cache.Unlock()
	return value.Value, true
}

func (cache *CacheMemory) Set(key string, value interface{}, ttl time.Duration) {
	cache.Lock()
	cache.Memory[key] = CacheElement{
		Value:       value,
		ErasureTime: time.Now().Add(ttl),
		TTL:         ttl,
	}
	cache.Unlock()
	//running a daemon to delete an item after the expiration date
	go func() {
		for {
			cache.Lock()
			if cache.Memory[key].ErasureTime.Unix() <= time.Now().Unix() {
				delete(cache.Memory, key)
				cache.Unlock()
				return
			}
			cache.Unlock()
			//to prevent the cache from being permanently blocked, we will run the daemon periodically
			time.Sleep(time.Second)
		}
	}()
}

func main() {
	var cache Cache = &CacheMemory{
		Memory: make(map[string]CacheElement),
	}
	//expiration test
	cache.Set("firstKey", "firstValue", 5*time.Second)
	time.Sleep(2 * time.Second)
	fmt.Println(cache.Get("firstKey"))
	time.Sleep(4 * time.Second)
	fmt.Println(cache.Get("firstKey"))
	time.Sleep(6 * time.Second)
	fmt.Println(cache.Get("firstKey"))
	//multiple access test
	for i := 0; i < 10; i++ {
		go func(i int) {
			cache.Set(string(i), i, 5*time.Second)
		}(i)
	}
	time.Sleep(time.Second)
	for i := 0; i < 3; i++ {
		go func(i int) {
			cache.Delete(string(i))
		}(i)
	}
	time.Sleep(time.Second)
	for i := 0; i < 11; i++ {
		go func(i int) {
			fmt.Println(cache.Get(string(i)))
		}(i)
	}
	time.Sleep(time.Second)
}
