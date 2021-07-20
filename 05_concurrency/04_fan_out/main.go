package main

import (
	"fmt"
	"sync"
	"time"
)

func generator() <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func worker(wg *sync.WaitGroup, name string, ch <-chan int) {
	defer wg.Done()

	for {
		time.Sleep(100 * time.Millisecond)

		i, ok := <-ch
		if ok {
			fmt.Printf("%s worker: %d\n", name, i)
		} else {
			fmt.Printf("%s worker is done\n", name)
			return
		}
	}
}

func main() {
	values := generator()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go worker(wg, "first", values)
	go worker(wg, "second", values)
	go worker(wg, "third", values)

	wg.Wait()
}
