package main

import (
	"fmt"
	"math"
	"os"
)

const overengineering = false

var r float64

// returns the volume of the sphere with radius r
func sphereVolume(r float64) (v float64, err error) {
	if r < 0 {
		err = fmt.Errorf("radius cannot be negative. Current value: %v", r)
		return
	}
	v = (4.0 / 3.0) * math.Pi * math.Pow(r, 3)
	return
}

func main() {
	if overengineering {
		fmt.Print("Enter radius of your sphere: ")

		_, err := fmt.Scanf("%v", &r)
		if err != nil {
			fmt.Printf("Your entered something wrong. Error: %v\n", err)
			os.Exit(1)
		}

		v, err := sphereVolume(r)
		if err != nil {
			fmt.Printf("Your entered something wrong. Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Value of the sphere: %v\n", v)
	} else {
		r = 12.3
		v, _ := sphereVolume(r)
		fmt.Printf("Value of the sphere with radius %v: %v\n", r, v)
	}
}
