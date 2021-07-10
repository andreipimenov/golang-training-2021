package notation

import (
	"github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/stack"
	"regexp"
	"strconv"
	"unicode"
)

type Lexem struct {
	String   string
	Number   float64
	IsNumber bool
}

func isOperator(c string) bool {
	if c == "+" || c == "-" || c == "*" || c == "/" {
		return true
	}
	return false
}

func precedence(c string) int {
	if c == "*" || c == "/" {
		return 1
	} else if c == "+" || c == "-" {
		return 0
	}
	return -1
}

func FromInfixToPostfix(infix string) []string {
	var lexems []string
	var numbers []string
	var postfix []string

	values := regexp.MustCompile("[-+*/()]").Split(infix, -1)

	for _, str := range values {
		if str != "" {
			numbers = append(numbers, str)
		}
	}

	flag := true
	count := 0

	for _, v := range infix {
		if unicode.IsDigit(v) {
			if flag {
				lexems = append(lexems, numbers[count])
				count++
				flag = false
			}
		} else {
			if char := string(v); isOperator(string(v)) || char == "(" || char == ")" {
				lexems = append(lexems, string(v))
				flag = true
			}
		}
	}

	//for _, v := range lexems {
	//	fmt.Println(v)
	//}

	var stack stack.StringStack

	for _, v := range lexems {
		if _, err := strconv.Atoi(v); err == nil {
			postfix = append(postfix, v)
		} else if _, err := strconv.ParseFloat(v, 8); err == nil {
			postfix = append(postfix, v)
		} else if v == "(" {
			stack.Push(v)
		} else if v == ")" {
			for {
				if top, err := stack.Top(); top == "(" || !err {
					break
				}
				temp, _ := stack.Pop()
				postfix = append(postfix, temp)
			}
			if top, err := stack.Top(); top == "(" && err {
				stack.Pop()
			}
		} else if isOperator(v) {
			if stack.IsEmpty() {
				stack.Push(v)
			} else {
				top, _ := stack.Top()
				if a, b := precedence(v), precedence(top); a > b {
					stack.Push(v)
				} else {
					for {
						top, _ = stack.Top()
						if el1, el2 := precedence(v), precedence(top); stack.IsEmpty() || el1 > el2 {
							break
						}
						postfix = append(postfix, top)
						stack.Pop()
					}
					stack.Push(v)
				}
			}
		}
	}

	for {
		pop, err := stack.Pop()
		if err == false {
			break
		}
		postfix = append(postfix, pop)
	}

	return postfix
}
