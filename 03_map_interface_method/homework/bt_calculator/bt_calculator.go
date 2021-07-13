/*
	Пакет для парсинга математических выражений и
	вычисления результата методом построения бинаного дерева
*/

package bt_calculator

import (
	"math"
	"strconv"
)

// Мапка операторов с их приоритетом для поиска точик перегиба
// https://habr.com/ru/post/263775/
var operators = map[string]int{
	"+": 3,
	"-": 3,
	"*": 2,
	"/": 2,
	"^": 1,
}

type Calculator struct {
}

// Функция должна построить дерево и вернуть слайс элементов
func (cal *Calculator) Calculate(expression string) float64 {
	// Если скобки открываются и закрываются - убираем
	expression = trimBrackets(expression)
	root := build(expression)
	return root.Calc()
}

func NewCalculator() *Calculator {
	return new(Calculator)
}

type Element struct {
	Data     string  // то, что передаем при создании элемента
	Value    float64 // float64 от даты, если дата - число
	IsLeaf   bool    // И если это число, то это - leaf
	Operator string  // + - * / ^
	Left     *Element
	Right    *Element
}

/*
	Поиск корня дерева. Это будет самый приоритетный оператор
	со следующими приоритетами

	"+" - 3

	"-" - 3

	"*" - 2

	"/" - 2

	"^" - 1

	Причем выбирается самый последний при одинаковом приоритете

	https://habr.com/ru/post/263775/
*/
func findInflectionPoint(in []byte) int {
	// слайс содержит массив приоритетов байт
	// не разделяю на лексемы, но т.к. ищем операторы, вполне ок
	// Если по индексу не оператор, там будет 0
	m := make([]int, len(in))
	// создадим магическое число, что уменьшит приоритет в скобках
	// если наткнулись на (, уменьшаем на -3, потом увеличиваем на 3 при )
	// Приоритет определяется как число из мапки operators + br
	br := 0
	for i, v := range in {
		if v == '(' {
			br -= 3
		}
		if v == ')' {
			br += 3
		}
		if p, ok := operators[string(v)]; ok {
			m[i] = p + br
		} else {
			m[i] = 0
		}
	}
	// идем справа налево и возвращаем индекс, по которому первый
	// попавшийся на пути максимальный приоритет
	var res, max int
	for i := len(in) - 1; i > 0; i-- {
		if m[i] > max {
			res = i
			max = m[i]
		}
	}
	return res
}

func trimBrackets(data string) string {
	// Режем скобки, сколько бы ни было
	in := []byte(data)
	for {
		if in[0] == '(' && in[len(in)-1] == ')' {
			in = in[1 : len(in)-1]
		} else {
			break
		}
	}
	return string(in)
}

func build(data string) *Element {
	// Основная функция. что строит дерево и врзвращает корень
	// вызывается рекурсивно
	e := new(Element)
	// Если скобки открываются и закрываются - убираем
	data = trimBrackets(data)
	//fmt.Println("building for", data)
	e.Data = data //Чисто для дебага, не актуально
	// а не число ли создаем
	v, err := strconv.ParseFloat(data, 64)
	if err == nil {
		//fmt.Println("get value", v)
		e.Value = v
		e.IsLeaf = true
		return e
	}
	// если не число, то конвертим в []byte, ищем точку перегиба
	in := []byte(data)
	inflectionIds := findInflectionPoint(in)
	// на точке перегиба - нужный нам оператор
	e.Operator = string(in[inflectionIds])
	//fmt.Println("oper", e.Operator)
	// ну и build'им левую и правую части
	leftData := string(in[:inflectionIds])
	//fmt.Println("leftData", leftData)
	if leftData != "" {
		e.Left = build(leftData)
	}
	rightData := string(in[inflectionIds+1:])
	//fmt.Println("rightData", leftData)
	if rightData != "" {
		e.Right = build(rightData)
	}
	return e
}

func (e *Element) Calc() float64 {
	// Максимально просто:
	// Если это лист - возвращаем число
	// Иначе считаем, вызывая Calc у детей
	if e.IsLeaf {
		return e.Value
	} else {
		switch e.Operator {
		case "+":
			return e.Left.Calc() + e.Right.Calc()
		case "-":
			return e.Left.Calc() - e.Right.Calc()
		case "*":
			return e.Left.Calc() * e.Right.Calc()
		case "/":
			return e.Left.Calc() / e.Right.Calc()
		case "^":
			return math.Pow(e.Left.Calc(), e.Right.Calc())
		default:
			// Сюда не должны дойти, но линтер ругается
			return e.Value
		}
	}
}
