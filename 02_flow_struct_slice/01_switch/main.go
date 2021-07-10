package main

import (
	"fmt"
)

func main() {
	x := 35
	// x = 15
	switch {
	case x > 10:
		fmt.Println("> 10")
		// fallthrough
	case x > 25:
		// break
		fmt.Println("> 25")
	default:
		fmt.Println("Default")
	}
}
