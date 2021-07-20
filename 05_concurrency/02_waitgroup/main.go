package main

import (
	"fmt"
	"sync"
)

func square(wg *sync.WaitGroup, i int) {
	fmt.Println(i * i)
	wg.Done()
}

func main() {
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	wg := sync.WaitGroup{}
	wg.Add(len(tasks))

	for _, task := range tasks {
		go square(&wg, task)
	}

	wg.Wait()
}
