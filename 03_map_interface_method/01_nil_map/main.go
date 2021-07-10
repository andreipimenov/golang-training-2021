package main

import "fmt"

func main() {
	var m map[string]int
	m = make(map[string]int)

	if m == nil {
		fmt.Println("Map is nil")
	} else {
		fmt.Println("Map is not nil")
	}

	m["Hello"] = 10

	value, ok := m["Hello"]
	fmt.Println(value, ok)
}
