package main

import (
	"flag"
	"fmt"
	"strconv"
)

//https://gist.github.com/bemasher/1777766
type Stack struct {
	top  *Element
	size int
}

type Element struct {
	value interface{}
	next  *Element
}

func (s *Stack) Len() int {
	return s.size
}

// Push a new element onto the stack
func (s *Stack) Push(value interface{}) {
	s.top = &Element{value, s.top}
	s.size++
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

//Взять элемент с вершины стека, не извлекая его
func (s *Stack) Top() (value interface{}) {
	if s.size > 0 {
		value = s.top.value
		return
	}
	return nil
}

//Описание алгоритма
//http://www.zrzahid.com/convert-to-reverse-polish-notation-and-evaluate-the-expression-shunting-yard-algorithm/
type Calc interface {
	Calculate(expression string) (float64, int)
}

type _calc_ struct{}

func (_calc_) Calculate(expression string) (float64, int) {
	literals := lex(expression)
	if literals == nil {
		fmt.Println("Error in function lex()")
		return 0, 1
	}
	literals = toReversePolishNotation(literals)
	if literals == nil {
		fmt.Println("Error in function toReversePolishNotation()")
		return 0, 1
	}
	result, err := eval(literals)
	if err != 0 {
		fmt.Println("Error in function eval()")
		return 0, 1
	}
	return result, 0
}

const (
	Digit byte = iota
	OpenBracket
	CloseBracket
	Div
	Mul
	Add
	Sub
	UnarySub
)

type Literal struct {
	LiteralType byte
	Value       float64
}

//Печать среза литералов
func printLiterals(literals []Literal) {
	for _, value := range literals {
		switch value.LiteralType {
		case Digit:
			fmt.Print("Digit:", value.Value, "\n")
		case OpenBracket:
			fmt.Print("OpenBracket:", value.Value, "\n")
		case CloseBracket:
			fmt.Print("CloseBracket:", value.Value, "\n")
		case Div:
			fmt.Print("Div:", value.Value, "\n")
		case Mul:
			fmt.Print("Mul:", value.Value, "\n")
		case Add:
			fmt.Print("Add:", value.Value, "\n")
		case Sub:
			fmt.Print("Sub:", value.Value, "\n")
		case UnarySub:
			fmt.Print("UnarySub:", value.Value, "\n")
		}
	}
}

//True, если литерал типа бинарный оператор
func isBinaryOperator(literal byte) bool {
	return literal == Div || literal == Mul || literal == Add || literal == Sub
}

//True, если символ цифра
func isDigit(sym byte) bool {
	return sym <= '9' && sym >= '0'
}

//Разбиение входной строки на срез литералов + упрощение унарных операций
func lex(input string) (literals []Literal) {
	for index := 0; index < len(input); index += 1 {
		var currentLiteral Literal
		var stringValue string
		if isDigit(input[index]) {
			currentLiteral.LiteralType = Digit
			stringValue += string(input[index])
			haveDot := false
			index += 1
			for index < len(input) {
				if input[index] == '.' {
					if haveDot == true {
						fmt.Println("Error: many dots in the operand")
						return nil
					}
					haveDot = true
					stringValue += string(input[index])
				} else if isDigit(input[index]) {
					stringValue += string(input[index])
				} else {
					index -= 1
					break
				}
				index += 1
			}
			value, err := strconv.ParseFloat(stringValue, 64)
			if err != nil {
				fmt.Println("Error: the number cannot be converted")
				return nil
			}
			currentLiteral.Value = value
		} else {
			switch input[index] {
			case '.':
				currentLiteral.LiteralType = Digit
				stringValue += string(input[index])
				index += 1
				for index < len(input) {
					if input[index] == '.' {
						fmt.Println("Error: many dots in the operand")
						return nil
					} else if isDigit(input[index]) {
						stringValue += string(input[index])
					} else {
						index -= 1
						break
					}
					index += 1
				}
				if stringValue == "." {
					stringValue = "0.0"
				}
				value, err := strconv.ParseFloat(stringValue, 64)
				if err != nil {
					fmt.Println("Error: the number cannot be converted")
					return nil
				}
				currentLiteral.Value = value
			case '(':
				currentLiteral.LiteralType = OpenBracket
			case ')':
				currentLiteral.LiteralType = CloseBracket
			case '/':
				currentLiteral.LiteralType = Div
			case '*':
				currentLiteral.LiteralType = Mul
			case '+':
				currentLiteral.LiteralType = Add
				if index == 0 {
					continue
				} else if literals[len(literals)-1].LiteralType == OpenBracket {
					continue
				} else if isBinaryOperator(literals[len(literals)-1].LiteralType) {
					continue
				}
			case '-':
				currentLiteral.LiteralType = Sub
				if index == 0 {
					currentLiteral.LiteralType = UnarySub
				} else if literals[len(literals)-1].LiteralType == OpenBracket {
					currentLiteral.LiteralType = UnarySub
				} else if isBinaryOperator(literals[len(literals)-1].LiteralType) {
					currentLiteral.LiteralType = UnarySub
				} else if literals[len(literals)-1].LiteralType == UnarySub {
					//-- = +
					literals = literals[:len(literals)-1]
					continue
				}
			case '=':
				return
			case ' ':
				continue
			default:
				fmt.Println("Error: unknown symbol")
				return nil
			}
		}
		literals = append(literals, currentLiteral)
	}
	return
}

//Вернет приоретет оператора
func getPriority(value Literal) int {
	if value.LiteralType == Add || value.LiteralType == Sub {
		return 2
	}
	if value.LiteralType == Mul || value.LiteralType == Div {
		return 3
	}
	return 4
}

//Приведение среза литералов к срезу литоралов в обратной польской нотации
func toReversePolishNotation(literals []Literal) (infix []Literal) {
	var stack Stack
	for _, value := range literals {
		if value.LiteralType == Digit {
			infix = append(infix, value)
		} else {
			if value.LiteralType == OpenBracket {
				stack.Push(value)
			} else if value.LiteralType == CloseBracket {
				for stack.Len() != 0 {
					topStack := stack.Top().(Literal)
					if topStack.LiteralType == OpenBracket {
						break
					}
					infix = append(infix, topStack)
					stack.Pop()
				}
				if stack.Len() == 0 {
					fmt.Println("Error: problems with brackets")
					return nil
				} else {
					stack.Pop()
				}
			} else {
				for stack.Len() != 0 {
					topStack := stack.Top().(Literal)
					if topStack.LiteralType == OpenBracket || getPriority(topStack) < getPriority(value) {
						break
					}
					infix = append(infix, topStack)
					stack.Pop()
				}
				stack.Push(value)
			}
		}
	}
	for stack.Len() != 0 {
		infix = append(infix, stack.Pop().(Literal))
	}
	return infix
}

//Функция вычисления выражения, приведенного к обратной польской нотации
func eval(literals []Literal) (float64, int) {
	var stack Stack
	for _, value := range literals {
		if isBinaryOperator(value.LiteralType) {
			if stack.Len() < 2 {
				fmt.Println("Error: Invalid expression")
				return 0, 1
			}
			op2 := stack.Pop().(Literal).Value
			op1 := stack.Pop().(Literal).Value
			var result Literal
			switch value.LiteralType {
			case Mul:
				result.Value = op1 * op2
			case Div:
				if op2 == 0.0 {
					fmt.Println("Error: division by zero")
					return 0, 1
				}
				result.Value = op1 / op2
			case Add:
				result.Value = op1 + op2
			case Sub:
				result.Value = op1 - op2
			}
			stack.Push(result)
		} else if value.LiteralType == UnarySub {
			if stack.Len() == 0 {
				fmt.Println("Error: Invalid expression")
				return 0, 1
			}
			op1 := stack.Pop().(Literal).Value
			var result Literal
			result.Value = op1 * (-1.0)
			stack.Push(result)
		} else {
			stack.Push(value)
		}
	}
	if stack.Len() != 1 {
		fmt.Println("Error: Invalid expression")
		return 0, 1
	}
	return stack.Pop().(Literal).Value, 0
}

func main() {
	calc := _calc_{}
	flag.Parse()
	if len(flag.Args()) > 0 {
		for _, each := range flag.Args() {
			result, err := calc.Calculate(each)
			if err == 0 {
				fmt.Println(result)
			}
		}
	} else {
		//написать свой пример
		result, err := calc.Calculate("20/2-(2+2*3)=")
		if err == 0 {
			fmt.Println(result)
		}
	}
}
