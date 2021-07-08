package main

import (
	"fmt"
	"strconv"
	"strings"
)

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

type Calc interface {
	Calculate() float64
}

type Expr string

func (e Expr) Validate() (bool, string) {
	allowed_symbols := "+-*/0123456789()"
	for _, char := range e {
		if !strings.Contains(allowed_symbols, string(char)) {
			return false, string(char)
		}
	}
	return true, ""
}

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
			i = j - 1
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
	fmt.Scan(&e)

	string_valid, ret_val := e.Validate()
	if string_valid {
		fmt.Println(e.Calculate())
	} else {
		fmt.Printf("Validation failed, character '%s' is not allowed. Shutting down.", ret_val)
	}
}