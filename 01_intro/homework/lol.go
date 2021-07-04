package main

import (
	"fmt"
	"math"
)


func main() {
	var r float64
	println("pls input radius")
	_, err := fmt.Scanf("%f",&r)
	if err != nil {
		return
	}
	q:=int(r*r*r*math.Pi)*4/3
	println(q)
}

