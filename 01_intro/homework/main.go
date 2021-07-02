//Homework by Vladimir glushakov

package main

import (
	"fmt"
	"math"
)

func vSphere(r float64) float64 {
	return 4.0 / 3.0 * math.Pi * math.Pow(r, 3)
}

func main() {
	r := 10.0
	fmt.Println(vSphere(r))
}
