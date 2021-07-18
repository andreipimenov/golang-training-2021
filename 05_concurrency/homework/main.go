package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/andreipimenov/golang-training-2021/05_concurrency/homework/gedis"
)

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

func main() {
	c := gedis.NewGedis()
	// example of simple set/get usage
	c.Set("a", 1, 10*time.Second)
	v, found := c.Get("a")
	fmt.Println("found?", found)
	if found {
		fmt.Println("a", v)
	}
	// key should be expired after 10 seconds (+ cleanup interval)
	time.Sleep(15 * time.Second)

	v, found = c.Get("a")
	fmt.Println("found?", found)
	if found {
		fmt.Println("a", v)
	}
	// how to delete key
	c.Set("b", 2, 10*time.Second)
	c.Delete("b")

	v, found = c.Get("b")
	fmt.Println("found?", found)
	if found {
		fmt.Println("b", v)
	}

	// lets test if expiry time is increasing
	c.Set("c", 3, 3*time.Second)
	go func() {
		for {
			c.Get("c")
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(5 * time.Second)
	v, found = c.Get("c")
	fmt.Println("found?", found)
	if found {
		fmt.Println("c", v)
	}

	// lets check concurrency
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			c.Set("go", i, 10*time.Second)
			if v, ok := c.Get("go"); ok {
				fmt.Printf("getting %v from goroutine#%v\n", v, i)
			}
			wg.Done()
		}(i, &wg)
	}
	wg.Wait()
}
