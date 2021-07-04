package main

import (
	"fmt"
	"math"
)

func sphereVolumeByRadius(r float64) float64 {
	return 4.0 / 3.0 * math.Pi * math.Pow(r, 3)
}

func main() {
	var r float64 = 2
	fmt.Printf("Volume of the sphere by radius %.0f is %.2f\n", r, sphereVolumeByRadius(r))
}
