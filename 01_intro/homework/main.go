package main

import (
	"fmt"
	"math"
)

var sphereRadius = 8.0

func sphereVolume(radius float64) float64 {
	return 4 * math.Pi * math.Pow(radius, 3) / 3
}

func main() {
	fmt.Println(sphereVolume(sphereRadius))
}
