package main

import (
	"fmt"
	"math"
	"os"
)

func main() {

	var radius float64

	fmt.Println("Enter the radius of a sphere to calculate its volume:")
	fmt.Fscan(os.Stdin, &radius)

	fmt.Println(volume(radius))
}

func volume(radius float64) float64 {
	return 4.0 / 3.0 * math.Pi * math.Pow(radius, 3)
}
