/*
	Adopted solution of my
	Test assignment for an internship for the Core Infrastructure team (Summer 2021)
	MisterZurg/vk-summer-internship-2021
*/
package rpn_calc

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// The task is implement the interface below
type Calc interface {
	Calculate(expression string) float64
}

func NewCalculator() *Calculator {
	return new(Calculator)
}

type StackRune []rune

func (s *StackRune) Push(v rune) {
	*s = append(*s, v)
}

func (s *StackRune) Pop() rune {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *StackRune) Peek() rune {
	return (*s)[len(*s)-1]
}

type StackFloat []float64

func (s *StackFloat) Push(v float64) {
	*s = append(*s, v)
}

func (s *StackFloat) Pop() float64 {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *StackFloat) Peek() float64 {
	return (*s)[len(*s)-1]
}

type Calculator struct{}

func (calc *Calculator) Calculate(expression string) float64 {
	expression = strings.ReplaceAll(expression, " ", "")
	isLegitimate(expression)
	preparedExpression := getAdoptedExpression(expression)
	//Решение подготовленного выражения в ОПН
	reversedPolishNotation := expression2ReversedPolishNotation(preparedExpression)
	arabicNumber := reversedPolishNotation2Answer(reversedPolishNotation)
	return arabicNumber
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

//Выражение в обратную польскую нотацию
func expression2ReversedPolishNotation(expression string) string {
	var current string
	var stack StackRune

	var currentPriority int
	// Проходимся по выражению посимвольно
	for i := 0; i < len(expression); i++ {
		// Получаем текущий приоритет
		currentPriority = getPriorityOfOperation(rune(expression[i]))
		// Если число
		if currentPriority == 0 {
			current += string(expression[i])
		}
		// Если открывающаяся скоба (
		if currentPriority == 1 {
			stack.Push(rune(expression[i]))
		}
		// Если математический операнд
		if currentPriority > 1 {
			// Разделяем элементы состоящие из более чем 1 цифры
			current += " "
			// Проверим стек на пустоту
			for len(stack) > 0 {
				/* Пока он не пустой, записываем в current
				   все элементы, с приоритетом меньше
				   текущего currentPriority
				*/
				if getPriorityOfOperation(rune(stack.Peek())) >= currentPriority {
					current += string(stack.Pop())
				} else {
					break
				}
			}
			stack.Push(rune(expression[i]))
		}
		// Если закрывающаяся скоба )
		if currentPriority == -1 {
			// Разделяем элементы состоящие из более чем 1 цифры
			current += " "
			// Забираем элементы из стека до тех пор,
			// пока не встретим открывающаяся скобку
			for getPriorityOfOperation(rune(stack.Peek())) != 1 {
				current += string(stack.Pop())
			}
			stack.Pop()
		}
	}
	for len(stack) > 0 {
		current += string(stack.Pop())
	}
	return current
}

//Обратная польская нотация в ответ, как ни странно :/
func reversedPolishNotation2Answer(rpn string) float64 {
	var operand string = ""
	var stack StackFloat

	for i := 0; i < len(rpn); i++ {
		// fmt.Println(string(rpn[i]))
		if (rune(rpn[i])) == ' ' {
			continue
		}
		// Если число
		if getPriorityOfOperation(rune(rpn[i])) == 0 {
			// fmt.Println("Число")
			// Собираем все число
			for rpn[i] != ' ' && getPriorityOfOperation(rune(rpn[i])) == 0 {
				operand += string(rpn[i])
				i++
				if i == len(rpn) {
					break
				}
			}

			// Для арабских чисел in case
			number, _ := strconv.ParseFloat(operand, 64)
			stack.Push(number)
			operand = ""
		}
		// Если математический операнд
		if getPriorityOfOperation(rune(rpn[i])) > 1 {
			// Забираем из стека два последних числа
			var a float64 = stack.Pop()
			var b float64 = stack.Pop()

			switch rune(rpn[i]) {
			case '+':
				stack.Push(b + a)
			case '-':
				stack.Push(b - a)
			case '*':
				stack.Push(b * a)
			case '/':
				if a == 0 {
					fmt.Println("error: division by zero!")
					panic("error: division by zero!")
				} else {
					ans := math.Floor(float64(b) / float64(a))
					// fmt.Println(math.Floor(float64(b) /float64(a)), "=" ,ans)
					stack.Push(ans)
				}
			}
		}
	}
	return stack.Pop()
}

func getPriorityOfOperation(token rune) int {
	if token == '*' || token == '/' {
		return 3
	} else if token == '+' || token == '-' {
		return 2
	} else if token == '(' {
		return 1
	} else if token == ')' {
		return -1
	} else {
		return 0 //Приоритет чисел
	}
}
