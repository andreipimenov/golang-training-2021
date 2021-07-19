package main

import (
	"fmt"
	"sync"
)

func producer(x int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			out <- x
		}
		close(out)
	}()
	return out
}

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	getAndPush := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go getAndPush(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	p1 := producer(1)
	p2 := producer(2)

	all := merge(p1, p2)
	for i := range all {
		fmt.Println(i)
	}
}
