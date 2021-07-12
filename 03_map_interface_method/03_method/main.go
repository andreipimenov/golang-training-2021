package main

import (
	"fmt"
)

type Temp float64

func (t Temp) Print() {
	fmt.Printf("%v\n", t)
}

// type Celcius Temp
type Celcius float64

func (c Celcius) PrintFormatted() {
	// c.Print()
	fmt.Printf("%.1f °C\n", c)
}

func main() {
	var t Temp = 28.0
	t.Print()

	var c Celcius = Celcius(t)
	c.PrintFormatted()
}
