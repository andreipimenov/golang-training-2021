package main

import (
	"fmt"
	"math"
)

func main() {
	var r float64
	fmt.Print("Enter the radius of the sphere: ")
	fmt.Scan(&r)
	result := (r * r * r * math.Pi * 4) / 3
	fmt.Println("The volume of the sphere:", result)
}
