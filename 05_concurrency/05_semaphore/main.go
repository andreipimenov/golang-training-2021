package main

import (
	"fmt"
	"time"
)

func worker(i int, limit <-chan struct{}) {
	fmt.Println(i)
	time.Sleep(2 * time.Second)
	<-limit
}

func main() {
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	limit := make(chan struct{}, 3)

	for _, task := range tasks {
		limit <- struct{}{}
		go worker(task, limit)
	}

	<-time.After(10 * time.Second)
}
