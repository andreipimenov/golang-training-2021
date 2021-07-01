package main

import "fmt"

const (
	x = 1 << 7
	// x = 1 << 62
	// x uint64 = 1 << 63
	// x = 1 << 200
	// y = x >> 199
)

var (
// z = 1 << 200
)

func main() {
	fmt.Printf("Binary: %b Type: %T Value: %v\n", x, x, x)
	// fmt.Printf("Binary: %b Type: %T Value: %v\n", y, y, y)
}
