package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Stack struct {
	top  *Element
	size int
}

type Element struct {
	value interface{}
	next  *Element
}

func (s *Stack) Empty() bool {
	return s.size == 0
}

func (s *Stack) Top() interface{} {
	return s.top.value
}

func (s *Stack) Push(value interface{}) {
	s.top = &Element{value, s.top}
	s.size++
}

func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

func Precedence(ch string) int {
	if ch == "+" || ch == "-" {
		return 1
	} else if ch == "*" || ch == "/" {
		return 2
	} else {
		return 0
	}
}

func HasHigherPrecedence(op1 string, op2 string) bool {
	op1Weight := Precedence(op1)
	op2Weight := Precedence(op2)
	return op1Weight >= op2Weight
}

func IsOperand(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func IsOperator(ch byte) bool {
	return strings.ContainsAny(string(ch), "+ & - & * & /")
}

func ToPostfix(s string) string {

	var stack Stack

	postfix := ""

	length := len(s)

	for i := 0; i < length; i++ {

		char := string(s[i])
		// Уберем пробелы
		if char == " " {
			continue
		}

		if char == "(" {
			stack.Push(char)
		} else if char == ")" {
			for !stack.Empty() {
				str, _ := stack.Top().(string)
				if str == "(" {
					break
				}
				postfix += " " + str
				stack.Pop()
			}
			stack.Pop()
		} else if !IsOperator(s[i]) {
			//не оператор => операнд
			j := i
			number := ""
			for ; j < length && IsOperand(s[j]); j++ {
				number = number + string(s[j])
			}
			postfix += " " + number
			i = j - 1
		} else {
			//символ оператор => попаем два элемента, выполняем операцию и пушим результат обратно
			for !stack.Empty() {
				top, _ := stack.Top().(string)
				if top == "(" || !HasHigherPrecedence(top, char) {
					break
				}
				postfix += " " + top
				stack.Pop()
			}
			stack.Push(char)
		}
	}

	for !stack.Empty() {
		str, _ := stack.Pop().(string)
		postfix += " " + str
	}

	return strings.TrimSpace(postfix)
}

func Operation(a int, b int, op string) int {
	switch op {
	case "+":
		return b + a
	case "-":
		return b - a
	case "*":
		return b * a
	case "/":
		return b / a
	default:
		return 0
	}
}

type Calc interface {
	Calculate() float64
}

type Expr string

func (e Expr) Calculate() float64 {
	e = Expr(ToPostfix(string(e)))
	var stack Stack
	for i := 0; i < len(e); i++ {
		char := e[i]
		if IsOperand(char) {
			j := i
			number := ""
			for ; j < len(e) && IsOperand(e[j]); j++ {
				number = number + string(e[j])
			}
			stack.Push(number)
		}
		switch char {
		case '+':
			b, _ := strconv.Atoi(stack.Pop().(string))
			a, _ := strconv.Atoi(stack.Pop().(string))
			stack.Push(strconv.Itoa(a + b))
		case '-':
			b, _ := strconv.Atoi(stack.Pop().(string))
			a, _ := strconv.Atoi(stack.Pop().(string))
			stack.Push(strconv.Itoa(a - b))
		case '*':
			b, _ := strconv.Atoi(stack.Pop().(string))
			a, _ := strconv.Atoi(stack.Pop().(string))
			stack.Push(strconv.Itoa(a * b))
		case '/':
			b, _ := strconv.Atoi(stack.Pop().(string))
			a, _ := strconv.Atoi(stack.Pop().(string))
			stack.Push(strconv.Itoa(a / b))
		}
	}

	ret_val, _ := strconv.Atoi(stack.Pop().(string))

	return float64(ret_val)
}

func main() {
	var e Expr
	e = "22*3/(5*4)+2"

	fmt.Println(e.Calculate())
}
