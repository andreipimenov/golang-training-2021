package main

import (
	"fmt"
	"math"
)

func sphereVolume(r int) float64 {
	return math.Pi * 4/3
}

func main() {
	var r int
	r = 6
	fmt.Printf("The volume of a sphere with radius %d is %f",r, sphereVolume(r))
}
