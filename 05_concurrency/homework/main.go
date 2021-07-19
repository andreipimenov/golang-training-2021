package main

import (
	"fmt"
	"github.com/andreipimenov/golang-training-2021/05_concurrency/homework/memorycache"

	"time"
)

func printInform(cache *memorycache.MemoryCache, key string) {
	if item, isHas := cache.Get(key); isHas {
		fmt.Printf("Cash have value for '%s' - ", key)
		fmt.Println(item)
	} else {
		fmt.Printf("Cash don't have value for '%s' !\n", key)
	}
}

func main() {
	cache := memorycache.InitializeMemoryCache()
	cache.Set("key one", "value one", time.Second)   // ttl for 'key one' 1 sec
	cache.Set("key two", "value two", 4*time.Second) // ttl for 'key two' 4 sec

	//cache.Delete("key one")

	printInform(cache, "key one")
	time.Sleep(time.Second * 4)
	printInform(cache, "key two")
}
