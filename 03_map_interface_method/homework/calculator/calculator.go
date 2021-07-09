package calculator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Calculator struct {
	input string
}

func (cal *Calculator) Calculate(expression string) float64 {
	tokens := splitInput([]byte(expression))
	rpn := reversePolishNotation(tokens)
	return calcRpn(rpn)
}

func NewCalculator() *Calculator {
	return &Calculator{}
}

type StringsStack struct {
	data []string
	Top  string
	Len  int
}

func (s *StringsStack) Push(e string) {
	s.data = append(s.data, e)
	s.Top = s.data[len(s.data)-1]
	s.Len = len(s.data)
}

func (s *StringsStack) Pop() (string, bool) {
	if len(s.data) == 1 {
		last := s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
		s.Len = len(s.data)
		s.Top = ""
		return last, true
	}
	if len(s.data) > 0 {
		last := s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
		s.Len = len(s.data)
		s.Top = s.data[len(s.data)-1]
		return last, true
	} else {
		return "", false
	}
}

func (s StringsStack) String() string {
	return strings.Join(s.data, " ")
}

// Разбиваем выражение на составные части: операнды и операторы (токены)
func splitInput(input []byte) (res []string) {
	// мапка операторов не содержит точку
	operators := make(map[byte]struct{})
	for _, v := range []byte("+-*/^()") {
		operators[v] = struct{}{}
	}
	// Создаем временный слайс слайсов байт
	tmpData := make([][]byte, 0)
	// Создаем слайс для текущего элмента tmpData
	curr := make([]byte, 0)
	for i, c := range input {
		if _, isOper := operators[c]; isOper {
			// Нашли оператор, скидываем в массив что накопили, если есть
			if len(curr) != 0 {
				tmpData = append(tmpData, curr)
			}
			// Обнуляем текущий элемент, аппендим оператор (сразу, без curr)
			curr = make([]byte, 0)
			tmpData = append(tmpData, []byte{c})
		} else {
			// число помещаем в curr, аппендим потом
			curr = append(curr, c)
		}
		// если дошли до конца, аппендим все, что есть (а есть последнее число)
		// а тут закрался баг на кейсе, когда в конце ")", аппендился []
		if i == len(input)-1 && len(curr) > 0 {
			tmpData = append(tmpData, curr)
		}
	}
	// делаем слайс стрингов - токенов, возвращаем его
	for _, val := range tmpData {
		res = append(res, string(val))
	}
	return
}

// Реализация алгоритма Обратной Польской записи
func reversePolishNotation(input []string) (res []string) {
	// Берем стэк
	stack := new(StringsStack)
	for idx, v := range input {
		// Итерируемся по каждому элеенту входного массива
		switch v {
		case "(":
			// Это всегда в стэк
			stack.Push(v)
		case ")":
			// Тут скидываем в результат все, что было в стэке до "("
			// Сами же скобки дропаем
			for stack.Len > 0 {
				if stack.Top != "(" {
					e, _ := stack.Pop()
					res = append(res, e)
				} else {
					_, ok := stack.Pop()
					if ok {
						break
					}
					break
				}
			}
		case "^":
			// Самая приоритетная операция, всегда в стэк
			//e, _ := stack.Pop()
			//res = append(res, e)
			stack.Push(v)
		case "*", "/":
			// убираем в стэк то, что не менее приоритетно, если есть
			// И кладем в стэк
			if stack.Len > 0 {
				t := stack.Top
				if t == "^" || t == "*" || t == "/" {
					e, _ := stack.Pop()
					res = append(res, e)
					stack.Push(v)
				} else {
					stack.Push(v)
				}
				// Если ничего не было, кладем в стэк
			} else {
				stack.Push(v)
			}
		case "+", "-":
			// Убираем все не менее приоритетное и кладем в стэк
			if stack.Len > 0 {
				t := stack.Top
				// возможно, я что-то делаю не так, но... ( - 2 ) ^ 2
				if input[idx-1] == "(" {
					res = append(res, "0")
				}
				// не менее приоритетными будут все операторы...
				if t == "^" || t == "*" || t == "/" || t == "+" || t == "-" {
					e, _ := stack.Pop()
					res = append(res, e)
					stack.Push(v)
				} else {
					stack.Push(v)
				}
			} else {
				stack.Push(v)
			}
		default:
			// операнды всегда в результат сразу
			res = append(res, v)
		}
		//fmt.Printf("v: %v\tstack: %v\tres:%v\n", v, stack, strings.Join(res, " "))
	}
	// Стэк перекладываем в результат
	for stack.Len > 0 {
		e, _ := stack.Pop()
		res = append(res, e)
	}
	return
}

// Считаем значение для выражения, данного в формате RPN (массив строк)
func calcRpn(rpn []string) float64 {
	// Мапка операторов
	operators := make(map[string]struct{})
	for _, v := range strings.Split("+ - * / ^", " ") {
		operators[v] = struct{}{}
	}
	// Делаем стэк
	stack := new(StringsStack)
	for _, v := range rpn {
		if _, ok := operators[v]; ok {
			// Нашли оператор. Достаем 2 значения со стека
			s1, _ := stack.Pop()
			s2, _ := stack.Pop()
			a, _ := strconv.ParseFloat(s1, 64)
			b, _ := strconv.ParseFloat(s2, 64)
			// И что-то считаем в зависимости от результата
			switch v {
			case "+":
				stack.Push(fmt.Sprintf("%f", a+b))
			case "-":
				stack.Push(fmt.Sprintf("%f", b-a))
			case "*":
				stack.Push(fmt.Sprintf("%f", a*b))
			case "/":
				stack.Push(fmt.Sprintf("%f", b/a))
			case "^":
				stack.Push(fmt.Sprintf("%f", math.Pow(b, a)))
			}
		} else {
			stack.Push(v)
		}
	}
	r, _ := stack.Pop()
	res, _ := strconv.ParseFloat(r, 64)
	return res
}
