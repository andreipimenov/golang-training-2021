package main

import (
	"fmt"
	"math"
)

func calculateSphereVolume(radius float64) float64 {
	var sphereVolume float64 = 4.0 / 3.0 * math.Pi * math.Pow(radius, 3)
	return sphereVolume

}

func main() {
	var sphereRadius float64 = 4.0
	fmt.Println(calculateSphereVolume(sphereRadius))

}
