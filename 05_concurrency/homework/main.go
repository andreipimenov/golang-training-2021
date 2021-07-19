// homework implements in-memory cache.
package main

import (
	"fmt"
	"sync"
	"time"
)

// Cache - interface that we have to implement
type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type MyInMemoryChache struct {
	sync.Mutex
	Items map[string]Item
	// DefaultExpirationTime time.Duration
}

type Item struct {
	Value interface{}
	TWC   time.Time
	TTL   time.Duration
}

func createItem(value interface{}, ttl time.Duration) Item {
	return Item{
		Value: value,
		TWC:   time.Now(),
		TTL:   ttl,
	}
}

// UpdateTimeSession - updates Items time of creation
func (i *Item) UpdateTimeSession() {
	i.TWC = time.Now()
}

// IsExpired - checks if Items TTL ge than current time
func (i *Item) IsExpired() bool {
	return i.TWC.Add(i.TTL).Unix() <= time.Now().Unix()
}

// Set - sets Item into our InMemoryChache implementation with given params
func (imc *MyInMemoryChache) Set(key string, value interface{}, ttl time.Duration) {
	imc.Lock()
	defer imc.Unlock()

	imc.Items[key] = createItem(value, ttl)
}

// Prevention panic: assignment to entry in nil map
func IntializeCache() *MyInMemoryChache {
	return &MyInMemoryChache{Items: make(map[string]Item)}
}

// Get - gets Item with given key
func (imc *MyInMemoryChache) Get(key string) (interface{}, bool) {
	imc.Lock()
	if item, ok := imc.Items[key]; ok {
		// Check if Item hasn't expired yet
		if !item.IsExpired() {
			item.UpdateTimeSession()
			imc.Unlock()
			return item.Value, true
		}
	}
	// If Item expired then we have to clear Chache
	imc.Unlock()
	imc.Delete(key)
	return nil, false
}

// Delete - removes Item from Items with given key
func (imc *MyInMemoryChache) Delete(key string) {
	imc.Lock()
	defer imc.Unlock()
	delete(imc.Items, key)
}

// TODO : Run It Somehow
func (imc *MyInMemoryChache) ClearExpiredItems(key string) {
	imc.Lock()
	for _, item := range imc.Items {
		if item.IsExpired() {
			imc.Unlock()
			imc.Delete(key)
		}
	}
	imc.Unlock()
}

func main() {

	myMemoryChache := IntializeCache()

	myMemoryChache.Set("ae2a6a4a-a4e1-485e-b3ab-b448b64fef46", "ae2a6a4a", time.Second)
	myMemoryChache.Set("bbe48869-c3bc-4c35-a670-b622f996b101", "bbe48869", time.Second*6)

	fmt.Println(myMemoryChache.Get("ae2a6a4a-a4e1-485e-b3ab-b448b64fef46"))
	fmt.Println(myMemoryChache.Get("bbe48869-c3bc-4c35-a670-b622f996b101"))
	fmt.Println(myMemoryChache.Get("notInMemoryChahce"))
	fmt.Println(myMemoryChache.Get("notInMemoryChahce2"))

	time.Sleep(time.Second * 5)

	fmt.Println(myMemoryChache.Get("ae2a6a4a-a4e1-485e-b3ab-b448b64fef46"))
	fmt.Println(myMemoryChache.Get("bbe48869-c3bc-4c35-a670-b622f996b101"))
	/*
		Worth mentioning:
		Assignment to entry in nil map
		fatal error: all goroutines are asleep - deadlock
		fatal error: sync: unlock of unlocked mutex
	*/
}
