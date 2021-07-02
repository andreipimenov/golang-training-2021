package main

import (
	"fmt"
	"math"
)

func sphereVolume(radius float64) (float64, error) {
	if radius < 0 {

		return 0, fmt.Errorf("math: negative radius value %v", radius)
	}

	return 4.0 / 3 * math.Pi * math.Pow(radius, 3), nil
}

func main() {
	var r float64 = 3
	res, err := sphereVolume(r)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

}
