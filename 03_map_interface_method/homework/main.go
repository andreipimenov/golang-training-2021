/*
	Task is to implement:
		type Calc interface {
		    Calculate(expression string) float64
		}
	Additionally:

	Add validation for your app to check if no symbols other than digits and allowed operators are passed as an argument to Calculate(string) method

	Use structs and methods, try to implement the task using tree-based approach and stack-based approach
*/
package main

import (
	"fmt"
	"github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/validation"

	// "github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/rpn_calc"
	"github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/bt_calc"
)

func main() {
	var promptedExpression string
	fmt.Scan(&promptedExpression)
	validMathExpression := validation.ValidateExpression(promptedExpression)
	calculator := bt_calc.NewCalculator()
	fmt.Println(calculator.Calculate(validMathExpression))
	fmt.Println(calculator.Calculate("20/2-(2+2*3)"))

	// calculatorRPNMethod := rpn_calc.NewCalculator()
	// fmt.Println(calculatorRPNMethod.Calculate("20/2-(2+2*3)"))
	// fmt.Println(calculatorRPNMethod.Calculate("20/2-(2+2*3)="))
	// fmt.Println(calculatorRPNMethod.Calculate("2 + (3 * 8) - (4 + (48 / (4 + 2)) * 6)"))
	// fmt.Println(calculatorRPNMethod.Calculate("((81 * 6) /42+ (3-1))"))
	// fmt.Println(calculatorRPNMethod.Calculate("-5/2"))

}
