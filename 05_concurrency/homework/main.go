package main

import (
	"fmt"
	"github.com/rabbit72/golang-training-2021/05_concurrency/homework/cache"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Second)
	myCache := cache.NewSimpleCache(1 * time.Second)
	defer myCache.Stop()
	myCache.Set("key1", "value1", 3*time.Second)
	myCache.Set("key2", "value3", 2*time.Second)

	for x := range ticker.C {
		_ = x
		fmt.Println(myCache)
		fmt.Println(myCache.Get("key2"))
	}
}
