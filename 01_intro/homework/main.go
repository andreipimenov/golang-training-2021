package main

import (
	"fmt"
	"math"
)

func sphereVolume(r float64) float64 {
	return math.Pow(r,3)*math.Pi * 4/3
}

func main() {
	var r float64
	r = 6
	fmt.Printf("The volume of a sphere with radius %f is %f",r, sphereVolume(r))
}
