package main

import "fmt"

// go run -race main.go

func main() {
	x := 0

	// mu := sync.Mutex{}

	for i := 0; i < 100; i++ {
		go func(i int) {
			//mu.Lock()
			x = i
			//mu.Unlock()
		}(i)
	}

	fmt.Println(x)
}
