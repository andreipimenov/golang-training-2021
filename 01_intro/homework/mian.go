package main

import (
	"fmt"
	"math"
)

func getSphereVolume(radius float64) float64 {
	return 4 * math.Pi * math.Pow(radius, 3.0) / 3
}

func main() {
	radius := 10.0
	fmt.Printf("Sphere Volume = %0.6f, R = %v\n", getSphereVolume(radius), radius)
}
