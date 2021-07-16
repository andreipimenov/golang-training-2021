package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println(i)
		}()
		// time.Sleep(2 * time.Second)
	}

	<-time.After(10 * time.Second)
}
