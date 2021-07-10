/*
	Пакет для парсинга математических выражений и
	вычисления результата методом приведения к обратной польской записи
*/

package rpn_calculator

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	stack "github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/stack"
)

// Нужна дл доставаний из стэка float64 вместо interface{}
var floatType = reflect.TypeOf(float64(0))

type Calculator struct{}

func (cal *Calculator) Calculate(expression string) float64 {
	tokens := splitInput([]byte(expression))
	rpn := reversePolishNotation(tokens)
	return calcRpn(rpn)
}

func NewCalculator() *Calculator {
	return &Calculator{}
}

// Функция преобразует interface{} из стэка в float64 для вычислений
func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
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
	Stack := stack.NewStack()
	for idx, v := range input {
		// Итерируемся по каждому элеенту входного массива
		switch v {
		case "(":
			// Это всегда в стэк
			Stack.Push(v)
		case ")":
			// Тут скидываем в результат все, что было в стэке до "("
			// Сами же скобки дропаем
			for Stack.Len > 0 {
				if Stack.Top.Value != "(" {
					e := Stack.Pop()
					res = append(res, fmt.Sprintf("%v", e))
				} else {
					Stack.Pop()
					break
				}
			}
		case "^":
			// Самая приоритетная операция, всегда в стэк
			//e, _ := Stack.Pop()
			//res = append(res, e)
			Stack.Push(v)
		case "*", "/":
			// убираем в стэк то, что не менее приоритетно, если есть
			// И кладем в стэк
			if Stack.Len > 0 {
				t := Stack.Top.Value
				if t == "^" || t == "*" || t == "/" {
					e := Stack.Pop()
					res = append(res, fmt.Sprintf("%v", e))
					Stack.Push(v)
				} else {
					Stack.Push(v)
				}
				// Если ничего не было, кладем в стэк
			} else {
				Stack.Push(v)
			}
		case "+", "-":
			// Убираем все не менее приоритетное и кладем в стэк
			if Stack.Len > 0 {
				t := Stack.Top.Value
				// возможно, я что-то делаю не так, но... ( - 2 ) ^ 2
				if input[idx-1] == "(" {
					res = append(res, "0")
				}
				// не менее приоритетными будут все операторы...
				if t == "^" || t == "*" || t == "/" || t == "+" || t == "-" {
					e := Stack.Pop()
					res = append(res, fmt.Sprintf("%v", e))
					Stack.Push(v)
				} else {
					Stack.Push(v)
				}
			} else {
				Stack.Push(v)
			}
		default:
			// операнды всегда в результат сразу
			res = append(res, v)
		}
		// дебажим тут https://www.semestr.online/informatics/polish.php
		//fmt.Printf("v: %v\tStack: %v\tres:%v\n", v, Stack, strings.Join(res, " "))
	}
	// Стэк перекладываем в результат
	for Stack.Len > 0 {
		e := Stack.Pop()
		res = append(res, fmt.Sprintf("%v", e))
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
	floatsStack := new(stack.Stack)
	for _, v := range rpn {
		if _, ok := operators[v]; ok {
			// Нашли оператор. Достаем 2 значения со стека
			a, _ := getFloat(floatsStack.Pop())
			b, _ := getFloat(floatsStack.Pop())
			// И что-то считаем в зависимости от результата
			switch v {
			case "+":
				floatsStack.Push(a + b)
			case "-":
				floatsStack.Push(b - a)
			case "*":
				floatsStack.Push(a * b)
			case "/":
				floatsStack.Push(b / a)
			case "^":
				floatsStack.Push(math.Pow(b, a))
			}
		} else {
			f, _ := strconv.ParseFloat(v, 64)
			floatsStack.Push(f)
		}
	}
	r, _ := getFloat(floatsStack.Pop())
	return r
}
