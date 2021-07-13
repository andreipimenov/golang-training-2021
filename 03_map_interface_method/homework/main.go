package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Calc interface {
	Calculate(expression string) float64
}

const (
	Number int = iota
	Operator
	LParen
	RParen
)

type Token struct {
	Type  int
	Value string
}

type Stack struct {
	Values []Token
}

func (s *Stack) Push(t Token) {
	s.Values = append(s.Values, t)
}

func (s *Stack) Pop() Token {
	l := s.Values[len(s.Values)-1]
	s.Values = s.Values[:len(s.Values)-1]
	return l
}

func (s *Stack) IsEmpty() bool {
	return len(s.Values) == 0

}

func (s *Stack) Top() Token {
	return s.Values[len(s.Values)-1]
}

func (s *Stack) Length() int {
	return len(s.Values)
}

type Calculator struct{}

var oprs = map[string]struct {
	prec int
	fn   func(x, y float64) float64
}{
	"+": {2, func(x, y float64) float64 { return x + y }},
	"-": {2, func(x, y float64) float64 { return x - y }},
	"*": {3, func(x, y float64) float64 { return x * y }},
	"/": {3, func(x, y float64) float64 { return x / y }},
}

func Parse(input string) Stack {
	var token string
	var tokens Stack
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "=", "")
	for _, v := range input {
		c := string(v)
		if unicode.IsDigit(v) || c == "." {
			token += c
		} else {
			if token != "" {
				tokens.Push(Token{Number, token})
				token = ""
			}
			if strings.ContainsAny(c, "+-*/") {
				tokens.Push(Token{Operator, c})
			} else if strings.Contains("()", c) {
				if c == "(" {
					tokens.Push(Token{LParen, c})
				} else {
					tokens.Push(Token{RParen, c})
				}
			}
		}
	}
	if token != "" {
		tokens.Push(Token{Number, token})
		token = ""
	}
	return tokens
}

func ConvertToRPN(input Stack) Stack {
	stack := Stack{}
	operators := Stack{}
	for _, v := range input.Values {
		switch {
		case v.Type == Operator:
			for !operators.IsEmpty() {
				if oprs[v.Value].prec <= oprs[operators.Top().Value].prec {
					stack.Push(operators.Pop())
					continue
				}
				break
			}
			operators.Push(v)

		case v.Type == LParen:
			operators.Push(v)
		case v.Type == RParen:
			for i := operators.Length() - 1; i >= 0; i-- {
				if operators.Values[i].Type != LParen {
					stack.Push(operators.Pop())
					continue
				} else {
					operators.Pop()
					break
				}
			}
		default:
			stack.Push(v)
		}
	}
	if !operators.IsEmpty() {
		for i := operators.Length() - 1; i >= 0; i-- {
			stack.Push(operators.Pop())
		}
	}
	return stack

}

func (c *Calculator) Calculate(expression string) float64 {
	input := Parse(expression)
	input = ConvertToRPN(input)
	stack := Stack{}
	for _, v := range input.Values {
		switch {
		case v.Type == Operator:
			f := oprs[v.Value].fn
			y, _ := strconv.ParseFloat(stack.Pop().Value, 64)
			x, _ := strconv.ParseFloat(stack.Pop().Value, 64)
			result := f(x, y)
			stack.Push(Token{Number, strconv.FormatFloat(result, 'f', 2, 64)})
		default:
			stack.Push(v)
		}

	}
	output, _ := strconv.ParseFloat(stack.Values[0].Value, 64)
	return output
}

func main() {
	var c Calc = &Calculator{}
	fmt.Println(c.Calculate("(5+6*(2-1) + 11)/2"))
}
