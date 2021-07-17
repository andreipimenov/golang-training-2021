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

type Object struct {
	value     interface{}
	timePoint time.Time     // Point of time when the user set value
	ttl       time.Duration // Time To Live
}

type CacheStruct struct {
	sync.Mutex
	object map[string]Object
}

func (c *CacheStruct) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.object, key)
}

func (c *CacheStruct) Get(key string) (interface{}, bool) {

	c.Lock()
	defer c.Unlock()

	valueItem, ok := c.object[key]
	if !ok {
		return nil, false
	}
	// update the object
	c.object[key] = Object{
		value:     valueItem.value,
		timePoint: time.Now(),
		ttl:       c.object[key].ttl,
	}
	return valueItem.value, true
}

func (c *CacheStruct) Set(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	c.object[key] = Object{
		value:     value,
		timePoint: time.Now(),
		ttl:       ttl,
	}
	// create goroutine for every value
	go func() {
		for {
			c.Lock()
			if c.object[key].timePoint.Add(c.object[key].ttl).Unix() <= time.Now().Unix() {
				delete(c.object, key)
				c.Unlock()
				return
			}
			c.Unlock()
		}
	}()
}

func main() {

	object := make(map[string]Object)

	cache := CacheStruct{
		object: object,
	}
	// sets the new value
	cache.Set("key", "Some Value", 4*time.Second)
	time.Sleep(2 * time.Second)
	// get value and update its TTL
	cache.Get("key")
	time.Sleep(3 * time.Second)
	// after sleep value still exists
	amount, ok := cache.Get("key")
	if ok {
		fmt.Println(amount)
	} else {
		fmt.Println("Nothing to show")
	}
	//  lets exceed TTl
	time.Sleep(5 * time.Second)
	amount, ok = cache.Get("key")
	if ok {
		fmt.Println(amount)
	} else {
		fmt.Println("Nothing to show")
	}
	// lets check delete function
	cache.Set("key", "Another Value", 4*time.Second)
	amount, ok = cache.Get("key")

	if ok {
		fmt.Println(amount)
	} else {
		fmt.Println("Nothing to show")
	}

	cache.Delete("key")
	amount, ok = cache.Get("key")

	if ok {
		fmt.Println(amount)
	} else {
		fmt.Println("Nothing to show")
	}
}
