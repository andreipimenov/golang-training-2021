package main

import "fmt"

func main() {
	x := 1
loop:
	for {
		switch {
		case x > 100:
			// break
			break loop
		default:
			x += x
		}
	}
	fmt.Println(x)
}
