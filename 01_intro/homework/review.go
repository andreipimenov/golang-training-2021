package main

import (
	"fmt"
)

func main() {
	// 1. .gitignore

	// 2. Variables
	vars()

	// 3. Types
	// types()
}

func vars() {
	{
		var r float64
		r = 10
		fmt.Printf("%T %.1f\n", r, r)
	}
	{
		var r float64 = 10
		fmt.Printf("%T %.1f\n", r, r)
	}
	{
		var r = 10.0
		// var r = 10
		fmt.Printf("%T %.1f\n", r, r)
	}
	{
		r := 10.0
		fmt.Printf("%T %.1f\n", r, r)
	}
	{
		r := float64(10)
		fmt.Printf("%T %.1f\n", r, r)
	}
}

func types() {
	{
		var x float64
		x = 5.0 / 3.0
		fmt.Printf("%T %v\n", x, x)
	}
	{
		var x float64
		x = 5 / 3
		fmt.Printf("%T %v\n", x, x)
	}
	{
		var x int
		x = 5 / 3
		fmt.Printf("%T %v\n", x, x)
	}
	{
		var x float64
		// var x int
		x = 5.0 / 3
		fmt.Printf("%T %v\n", x, x)
	}
	{
		// var x float64
		var x int
		var five int = 5
		var three int = 3
		x = five / three
		fmt.Printf("%T %v\n", x, x)
	}

}
