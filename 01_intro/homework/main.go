package main

import (
	"fmt"
	"math"
)

func calculateSphereVolume(radius float64) float64 {
	radius = math.Max(radius, 0.0)
	return 4.0 / 3.0 * math.Pi * math.Pow(radius, 3)
}

func main() {
	radius := 5.0
	fmt.Printf(
		"Volume of the sphere of the radius %v is %v\n",
		radius,
		calculateSphereVolume(radius),
	)
}
