package main

import (
	"fmt"
	"log"
	"math"
)

const PI = 3.14159

// V = 4/3 * πr ^ 3
func volumeOfTheSphere(r float64) (float64, error) {
	if r < 0 {
		return 0, fmt.Errorf("You've entered negative radius %v", r)
	}
	return 4 / 3 * PI * math.Pow(r, 3), nil
}

func main() {
	var radius float64
	fmt.Println("Enter radius of sphere: ")
	_, err := fmt.Scan(&radius)
	if err != nil {
		log.Fatal(err)
	} else {
		answer, err := volumeOfTheSphere(radius)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(answer)
	}
}
