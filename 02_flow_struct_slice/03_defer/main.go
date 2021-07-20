package main

import "fmt"

func main() {
	i := 5
	defer fmt.Println("Deferred", i)
	i = i * 2
	fmt.Println(i)
}
