package main

import (
	"fmt"

	"github.com/andreipimenov/golang-training-math/v2/math"
)

func main() {
	fmt.Println(math.Sum(2, 3))
	fmt.Println(math.Multiply(2, 3))

	fmt.Println(math.Sum(1, 2, 3, 4, 5))
}
