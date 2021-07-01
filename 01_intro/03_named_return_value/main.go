package main

import "fmt"

func split(sum int) (x, y int) {
	// change to := to test if x is initialized
	x = sum * 4 / 9
	y = sum - x
	return
}

func main() {
	fmt.Println(split(17))
}
