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
	value       interface{}
	timeCreated time.Time
	ttl         time.Duration
}

type memoryCache struct {
	sync.Mutex
	items map[string]Item
}

func (it *memoryCache) Set(key string, value interface{}, ttl time.Duration) {
	it.Lock()
	var val = Item{
		value:       value,
		timeCreated: time.Now(),
		ttl:         ttl,
	}
	it.items[key] = val
	it.Unlock()
	go func() {
		for {
			it.Lock()
			if time.Now().UnixNano()-it.items[key].timeCreated.UnixNano() >
				it.items[key].ttl.Nanoseconds() {
				it.Unlock()
				it.Delete(key)
				return
			}
			it.Unlock()
		}
	}()
}

func (it *memoryCache) Get(key string) (interface{}, bool) {
	it.Lock()
	if _, ok := it.items[key]; !ok {
		it.Unlock()
		return nil, false
	} else {
		it.Unlock()
		return it.items[key].value, true
	}
}

func (it *memoryCache) Delete(key string) {
	it.Lock()
	defer it.Unlock()
	delete(it.items, key)
}
func main() {
	items := make(map[string]Item)
	fmt.Println(items)
	var mem memoryCache = memoryCache{
		items: items,
	}
	//Simple tests
	mem.Set("pack", "5051", 2*time.Second)
	time.Sleep(1 * time.Second)
	fmt.Println(mem.Get("pack"))
	mem.Set("crack", "05", 5*time.Second)
	time.Sleep(2 * time.Second)
	fmt.Println(mem.Get("crack"))
	time.Sleep(1 * time.Second)
	mem.Delete("crack")
	fmt.Println(mem.Get("crack"))

}
