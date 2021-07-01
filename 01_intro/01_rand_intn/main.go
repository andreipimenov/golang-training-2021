package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// uncomment to randomize
	// rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		fmt.Println("My favorite number is", rand.Intn(10))
	}
}
