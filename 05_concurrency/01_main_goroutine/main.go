package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func noSleep(s string) {
	fmt.Println(s)
}

func main() {
	go say("world")
	noSleep("hello")

	// time.Sleep(1 * time.Second)
}
