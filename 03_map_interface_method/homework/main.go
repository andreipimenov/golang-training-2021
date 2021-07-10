package main

import (
	"fmt"
	"github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/notation"
	"github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/stack"
	"strconv"
)

var operations = map[string]func(a float64, b float64) float64{
	"+": func(a float64, b float64) float64 { return a + b },
	"-": func(a float64, b float64) float64 { return a - b },
	"*": func(a float64, b float64) float64 { return a * b },
	"/": func(a float64, b float64) float64 { return a / b },
}

type Calculator struct{}

type Calc interface {
	Calculate(expression string) float64
}

func (c Calculator) Calculate(expression string) float64 {
	postfix := notation.FromInfixToPostfix(expression)

	var stack stack.FloatStack

	for _, v := range postfix {
		if num, err := strconv.Atoi(v); err == nil {
			stack.Push(float64(num))
		} else if num, err := strconv.ParseFloat(v, 8); err == nil {
			stack.Push(num)
		} else {
			if v != " " {
				a, e1 := stack.Pop()
				b, e2 := stack.Pop()

				if e1 == true && e2 == true {
					val := operations[v](b, a)
					stack.Push(val)
				}
			}
		}
	}

	pop, _ := stack.Pop()

	return pop
}

func main() {

	expression := "20/2-(2+2*3)" // 2
	//expression := "17+5*(4+8-7)*89/(5-9*8)" // −16,2089552239
	//expression := "((17+5*(4+(8-7))*89)/(5-9*8))" // −33,4626865672
	//expression := "121/(17/5-7)/12" // −2,80092592593
	//expression := "5*85.5" // 427.5
	//expression := "147*((58+45)/(78/45))*124/(121+56*8)" // 1903,62714614
	//expression := "(121/(17.5/5-7)/12.4)*(148-12.5)" // −377,776497696
	//expression := "(121/(17/(5-7/12)))+(148-12)" // 167,43627451
	//expression := "(121/(17/(5-7/12)))/(148/(12)*(121/(17.5/5-7)/12.4)*(148)-12.5)" // −0,00616208691

	// Don't support unary minus !
	// Support float value (like 1.01, not like 1,01)

	fmt.Print(Calculator{}.Calculate(expression))
}
