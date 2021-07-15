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
	Value interface{}
	Expiration int64
	TTL time.Duration
}

type MemCache struct {
	sync.Mutex
	defaultExp time.Duration
	cleanupInterval time.Duration
	items map[string]Item
}

func New(defaultExp, cleanupInterval time.Duration) *MemCache {
	//func to create a new cache with a default expiration timer and a cleanupInterval for the GC implementation
	items := make(map[string]Item)

	cache := MemCache{
		items: items,
		defaultExp: defaultExp,
		cleanupInterval: cleanupInterval,
	}

	//GC is started with a given cleanupInterval
	if cleanupInterval > 0 {
		cache.StartCleanUp()
	}

	return &cache
}

func (c *MemCache) StartCleanUp() {
	//each cache container has its own GC running concurrently
	go c.CleanUp()
}

func (c *MemCache) CleanUp() {
	for {
		//wait for the cleanupInterval to pass
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.getExpiredKeys(); len(keys) != 0 {
			c.Lock()
			for _, key := range keys {
				delete(c.items, key)
			}
			c.Unlock()
		}

	}
}

func (c *MemCache) getExpiredKeys() (keys []string) {
	c.Lock()
	defer c.Unlock()

	//find expired elements and delete them from the map
	for key, item := range c.items {
		if time.Now().UnixNano() > item.Expiration && item.Expiration > 0 {
			keys = append(keys, key)
		}
	}

	return
}

func (c *MemCache) Set(key string, value interface{}, ttl time.Duration) {
	var expiration int64

	if ttl <= 0 {
		ttl = c.defaultExp
	}

	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = Item{
		Value: value,
		Expiration: expiration,
		TTL: ttl,
	}
}

func (c *MemCache) Get(key string) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()
	
	cValue, found := c.items[key]

	if !found {
		return nil, false
	}

	if cValue.Expiration > 0 {
		if time.Now().UnixNano() > cValue.Expiration {
			return nil, false
		}
	}

	//renew TTL after accessing the item
	cValue.Expiration = time.Now().Add(cValue.TTL).UnixNano()

	return cValue.Value, true

}

func (c *MemCache) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.items, key)
}

func main() {

	//Low time is used only for example, increase the cleanupInterval for actual use, since
	//frequent GC calls decrease performance.
	cache := New(time.Second, 3*time.Second)

	cache.Set("key1", "NeededValue", 2*time.Second)

	item, found := cache.Get("key1")

	if found {
		fmt.Println(item)
	} else {
		fmt.Println("item not found")
	} //Prints the value for item in cache.

	//Wait long enough for the GC to kick in.
	time.Sleep(5*time.Second)

	item, found = cache.Get("key1")
	//Since we waited for 5 seconds and TTL of the item was 2 seconds the GC deleted it and we shouldn't be able
	//to access it.
	if found {
		fmt.Println(item)
	} else {
		fmt.Println("item not found")
	} //Prints "item not found"

}
