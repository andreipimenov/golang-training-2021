// TO DO
// Добавить поддержку отрицательных значений
// Добавить проверку пользовательского ввода
// Вынести стек отдельно

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"unicode"
)

type NumStack []float64

func (s *NumStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *NumStack) Push(str float64) {
	*s = append(*s, str)
}

func (s *NumStack) Pop() (float64, bool) {
	if s.IsEmpty() {
		return 0, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

type ValuesStack []string

func (s *ValuesStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *ValuesStack) Top() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		return (*s)[len(*s)-1], true
	}
}

func (s *ValuesStack) Push(str string) {
	*s = append(*s, str)
}

func (s *ValuesStack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

type Calculator struct{}

type Calc interface {
	Calculate(expression string) float64
}

func isOperator(c string) bool {
	switch c {
	case "+": return true
	case "-": return true
	case "*": return true
	case "/": return true
	default: return false
	}
}

func operationOrder(c string) int { // присваиваем приоритетность операциям для работы со стеком
	switch c {
	case "+": return 0
	case "-": return 0
	case "*": return 1
	case "/": return 1
	default: return -1
	}
}

func expressionSort(defaultExpression string) []string {

	var numbers []string
	var parsedValues []string
	var sortedArray []string

	values := regexp.MustCompile("[-+*/()]").Split(defaultExpression, -1)

	for _, num := range values {
		if num != "" {
			numbers = append(numbers, num)
		}
	}

	isNewNumber := true
	counter := 0

	for _, value := range defaultExpression {
		if unicode.IsDigit(value) {
			if isNewNumber && numbers != nil {
				parsedValues = append(parsedValues, numbers[counter]) // подтягиваем значение из среза numbers
				counter++
				isNewNumber = false // не рассматриваем числа до следующего знака
			}
		} else {
			if char := string(value); isOperator(string(value)) || char == "(" || char == ")" {
				parsedValues = append(parsedValues, string(value))
				isNewNumber = true
			}
		}
	}

	var parserStack ValuesStack

	for _, value := range parsedValues {
		if _, err := strconv.Atoi(value); err == nil {
			sortedArray = append(sortedArray, value)
		} else if _, err := strconv.ParseFloat(value, 8); err == nil {
			sortedArray = append(sortedArray, value)
		} else if value == "(" {
			parserStack.Push(value)
		} else if value == ")" {
			for {
				if top, err := parserStack.Top(); top == "(" || !err {
					break
				}
				temp, _ := parserStack.Pop()
				sortedArray = append(sortedArray, temp)
			}
			if top, err := parserStack.Top(); top == "(" && err {
				parserStack.Pop()
			}
		} else if isOperator(value) {
			if parserStack.IsEmpty() {
				parserStack.Push(value)
			} else {
				top, _ := parserStack.Top()
				if a, b := operationOrder(value), operationOrder(top); a > b { // анализируем приоритетность операций в стеке
					parserStack.Push(value)
				} else {
					for {
						top, _ = parserStack.Top()
						if el1, el2 := operationOrder(value), operationOrder(top); parserStack.IsEmpty() || el1 > el2 {
							break
						}
						sortedArray = append(sortedArray, top)
						parserStack.Pop()
					}
					parserStack.Push(value)
				}
			}
		}
	}

	for {
		pop, err := parserStack.Pop()
		if err == false {
			break
		}
		sortedArray = append(sortedArray, pop)
	}

	return sortedArray
}

func (c Calculator) Calculate(expression string) float64 {

	sortedExpression := expressionSort(expression)

	var calcStack NumStack

	for _, v := range sortedExpression {
		if num, err := strconv.Atoi(v); err == nil {
			calcStack.Push(float64(num))
		} else if num, err := strconv.ParseFloat(v, 8); err == nil {
			calcStack.Push(num)
		} else {
			if v != " " {
				b, lastValue := calcStack.Pop()
				a, preLastValue := calcStack.Pop()

				if lastValue == true && preLastValue == true {
					switch v {
					case "+": calcStack.Push(a+b)
					case "-": calcStack.Push(a-b)
					case "*": calcStack.Push(a*b)
					case "/": calcStack.Push(a/b)
					}
				}
			}
		}
	}

	pop, _ := calcStack.Pop()

	return pop
}

func main() {

	fmt.Println("Enter an expression without whitespaces using only +,-,/,*,(,) and positive numbers:")

	var expression string

	fmt.Fscan(os.Stdin, &expression)

	fmt.Println("Result:")

	fmt.Print(Calculator{}.Calculate(expression))
}
