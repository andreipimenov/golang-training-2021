package main

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	value interface{}
	ttl   time.Duration
}

type CacheBuffer struct {
	sync.Mutex
	items    map[string]*Item
	cacheTTL map[string]time.Duration //
}

//worker sets off timer before key invalidation, after key is invalid it gets deleted and its ttl from ttl map
func (cb *CacheBuffer) worker(key string) {
	cb.Lock()
	ttl := cb.items[key].ttl
	cb.Unlock()
	<-time.After(ttl)
	cb.Lock()
	delete(cb.items, key)
	delete(cb.cacheTTL, key)
	cb.Unlock()
	fmt.Println("Deleted item with key: ", key)
}

func (cb *CacheBuffer) Set(key string, value interface{}, ttl time.Duration) {
	cb.Lock()
	cb.items[key] = &Item{value: value, ttl: ttl}
	cb.cacheTTL[key] = ttl
	go cb.worker(key)
	cb.Unlock()
}

func (cb *CacheBuffer) Get(key string) (interface{}, bool) {
	_, found := cb.items[key]
	if found {
		cb.Lock()
		temp := cb.items[key].value
		cb.items[key].ttl += cb.cacheTTL[key]
		cb.Unlock()
		return temp, true
	} else {
		return nil, false
	}
}

func (cb *CacheBuffer) Delete(key string) {
	_, found := cb.items[key]
	if found {
		cb.Lock()
		cb.items[key].ttl = 0
		cb.Unlock()
	}
}

func (cb *CacheBuffer) DisplayItems() {
	cb.Lock()
	fmt.Println(cb.items)
	cb.Unlock()
}

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
	DisplayItems()
}

func main() {
	var memoryCache Cache = &CacheBuffer{
		items:    make(map[string]*Item),
		cacheTTL: make(map[string]time.Duration),
	}
	memoryCache.Set("example1", 1, 3*time.Second)
	memoryCache.Set("example2", 2, 2*time.Second)
	memoryCache.Set("example3", 3, 1*time.Second)
	time.Sleep(2 * time.Second) //when time Expires all three items will be deleted in order because of ttl
	//memoryCache.Set("example1", 1, 2*time.Second)
	//memoryCache.Set("example2", 2, 2*time.Second)
	//memoryCache.Set("example3", 3, 1*time.Second)
	//memoryCache.Get("example1")
	//memoryCache.Get("example1")
	//time.Sleep(2 * time.Second)
	////message of deletion of example1 won't be printed because after Get its time of existence will be increased by
	////its ttl
	//memoryCache.Set("example1", 1, 2*time.Second)
	//memoryCache.Set("example2", 2, 2*time.Second)
	//memoryCache.Set("example3", 3, 1*time.Second)
	//memoryCache.Get("example1")
	//memoryCache.Get("example1")
	//time.Sleep(6 * time.Second)
	////all messages of deletion of example1 will be printed because sleep is the same as the longes time of existence
	////of example1 6 seconds
	//memoryCache.Set("example1", 1, 2*time.Second)
	//memoryCache.Set("example2", 2, 2*time.Second)
	//memoryCache.Set("example3", 3, 1*time.Second)
	//memoryCache.Get("example1")
	//memoryCache.Get("example1")
	//memoryCache.Delete("example1")
	//time.Sleep(6 * time.Second)
	////example one will be deleted first because Delete sets time of existence of item to zero
}
