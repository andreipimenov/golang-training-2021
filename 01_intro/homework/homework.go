package main

import (
	"fmt"
	"math"
)

func sphereVolumeCalc(radius float64) float64 {
	return 4.0 / 3.0 * math.Pi * math.Pow(radius, 3)
}

func main() {
	var radius float64

	fmt.Printf("Please enter the radius of the sphere: ")
	fmt.Scan(&radius)
	fmt.Println("Sphere volume is:", sphereVolumeCalc(radius))
}