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

type CacheMap struct {
	m       map[string]Value
	counter int
	mu      sync.Mutex
}

type Value struct {
	value       interface{}
	ttl         time.Duration
	currentTime time.Time
}

func (t *CacheMap) Set(key string, value interface{}, ttl time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var valueFromMap = Value{value, ttl, time.Now()}
	t.m[key] = valueFromMap
	t.counter++
	if t.counter > 1000000 {
		t.clearCache()
		t.counter = 0
	}

}

func (t CacheMap) Get(key string) (interface{}, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if k, ok := t.m[key]; ok {
		if time.Now().Sub(k.currentTime) < k.ttl {
			return k.value, ok
		}

	}
	return "", false
}

func (t CacheMap) Delete(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.m, key)
}

//function, clearing cache after million requests

func (t CacheMap) clearCache() {
	fmt.Println(t.m)
	for key, element := range t.m {
		if time.Now().Sub(element.currentTime) > element.ttl {
			delete(t.m, key)
		}

	}

}

func main() {
	var m1 map[string]Value
	m1 = make(map[string]Value)

	var cache CacheMap = CacheMap{m: m1, counter: 0}

	cache.Set("cache1", "123", time.Second)
	time.Sleep(2 * time.Second)
	cache.Set("cache2", "456", 4*time.Second)
	fmt.Println(cache.Get("cache1"))
	fmt.Println(cache.Get("cache2"))
	cache.Set("cache1", "123", time.Second)
	cache.Delete("cache1")
	//cache.clearCache()

}
