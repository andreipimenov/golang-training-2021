package main

import (
	"fmt"
)

func zero() int {
	return 0
}

func division() {
	v := 3 / zero()
	fmt.Println(v)
}

func main() {
	v := 3 / zero()
	fmt.Println(v)

	// division()

	// go division()
	// time.Sleep(time.Second)

	// panic("Hey")

	// r := recover()
	// fmt.Println(r)

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("Deferred: %v\n", r)
	// 		return
	// 	}
	// 	fmt.Println("No panics")
	// }()
}
