package main

import (
	"fmt"
	"math"
)

func calculateSphereVolume(diameter float64) float64 {
	return math.Pi * math.Pow(diameter, 3.0) / 6
}

func main() {
	const diameter = 6.00
	fmt.Printf("Our Sphere = %0.3f, Diameter = %v\n", calculateSphereVolume(diameter), diameter)

}
