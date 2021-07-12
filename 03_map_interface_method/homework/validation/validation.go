package validation

import (
	"fmt"
	"strings"
)

func ValidateExpression(expression string) string {
	isLegitimate(expression)
	return getAdoptedExpression(expression)
}

func isLegitimate(expression string) {
	numberOfLeftBrackets := strings.Count(expression, "(")
	numberOfRightbrackets := strings.Count(expression, ")")
	if numberOfLeftBrackets != numberOfRightbrackets {
		fmt.Println("error: number of left brackets not equal number of right")
		panic("error: number of left brackets not equal number of right")
	}

	// Check expressions like (I+II)*
	if expression[0] == '*' || expression[len(expression)-1] == '*' ||
		expression[len(expression)-1] == '-' || // ха-ха  expression[0] == '-'
		expression[0] == '+' || expression[len(expression)-1] == '+' ||
		expression[0] == '/' || expression[len(expression)-1] == '/' {
		fmt.Println("error: beginning/ending of expression contains illegal symbol")
		panic("error: beginning/ending of expression contains illegal symbol")
	}
}

func getAdoptedExpression(expression string) string {
	if expression[len(expression)-1] == '=' {
		expression = expression[:len(expression)-1]
	}

	expression = strings.ReplaceAll(expression, " ", "")

	var adoptedExpression string = ""
	for token := 0; token < len(expression); token++ {
		var symbol = rune(expression[token])
		if symbol == '-' {
			if token == 0 { // Чекаем первый символ, явл ли он -
				adoptedExpression += "0"
			} else if expression[token-1] == '(' {
				adoptedExpression += "0"
			}
		}
		adoptedExpression += string(symbol)
	}
	return adoptedExpression
}
