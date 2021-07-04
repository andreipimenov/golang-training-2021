package main

import (
	"fmt"
	"math"
)

func calculateSphereVolume(r float64) float64 {
	if r <= 0 {
		return 0
        }
	return 4.0 / 3.0  * math.Pi * math.Pow(r, 3)
}

func main() {
	r := 5.0
	fmt.Println(calculateSphereVolume(r))
}
