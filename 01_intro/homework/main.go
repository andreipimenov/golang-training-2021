package main

import (
	"fmt"
	"math"
)

func main() {
	const PI = math.Pi
	var r float64 = 5
	fmt.Println("The volume of the sphere is", sphereVol(r, PI))
}

func sphereVol(r, PI float64) float64 {
	return (4 * math.Pow(r, 3) / 3) * PI
}
