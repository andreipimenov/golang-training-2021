package main

import (
	"fmt"
	"strings"
)

var in = "20/2-(2+2*3)^2="

func genArrowHelp(pos int) string {
	res := make([]rune, pos+1)
	for i := 0; i < pos; i++ {
		res[i] = ' '
	}
	res[pos] = '^'
	return string(res)
}

func inputValidate(in string) (err error) {
	// 0. Заменяем пробелы, убираем знак равенства в конце, если есть
	input := []byte(strings.Replace(in, " ", "", -1))
	if input[len(input)-1] == '=' {
		input = input[0 : len(input)-1]
	}
	// 1 Чекаем наличие только символов [0-9], + - * / ( ) ^ . =.
	// самый изи способ - регулярка
	// matched, err := regexp.Match("[^0-9-+/*^.()=]+", []byte(in))
	// но используется сет для проверки, чтобы указать на ошибку
	// Сначала делаем мапку, по ней будем чекать
	validChars := make(map[byte]struct{})
	for _, v := range []byte("0123456789+-*/()^.") {
		validChars[v] = struct{}{}
	}
	// проверяем и возвращаем ошибку, если символ не нашелся в валидных
	for idx, c := range input {
		if _, ok := validChars[c]; !ok {
			return fmt.Errorf("%v: invalid character '%v' on position %v\n%v",
				in, string(c), idx, genArrowHelp(idx))
		}
	}
	// 2. Проверяем скобки
	// 2.1 Проверка равного количества открывающихся и закрывающихся скобок
	// 2.2 проверка, что они сначала открываются и просто не закрываются
	// 2.3 проверка, что они не идут подряд
	opening, closing, lastOpening := 0, 0, 0
	for idx, c := range input {
		switch c {
		case '(':
			opening += 1
			lastOpening = idx
		case ')':
			if idx-1 == lastOpening {
				return fmt.Errorf("%v: closing bracket right after the opening on position %v\n%v",
					in, idx, genArrowHelp(idx))
			}
			closing += 1
			if closing > opening {
				return fmt.Errorf("%v: closing bracket before the opening on position %v\n%v",
					in, idx, genArrowHelp(idx))
			}
		}
	}
	if opening > closing {
		return fmt.Errorf("%v: unclosed bracket detected on position %v\n%v",
			in, lastOpening, genArrowHelp(lastOpening))
	}
	// 3. Проверка операторов
	// 3.1 Проверяем, что арифметические операторы не идут подряд
	// 3.2 Проверяем, что выражение не заканчивается оператором
	// 3.3 А оператор вообще есть?
	// 3.4 В начале может быть только "-"
	operators := make(map[byte]struct{})
	for _, v := range []byte("+-*/^.") {
		operators[v] = struct{}{}
	}
	for idx, c := range input {
		if _, isOper := operators[c]; isOper {
			if idx == 0 && c != '-' {
				return fmt.Errorf("%v: not '-' operator on the 1st  position\n%v",
					in, genArrowHelp(idx))
			}
			if idx == len(input)-1 {
				return fmt.Errorf("%v: operator on the last position\n%v",
					in, genArrowHelp(idx))
			}
			if _, isNextOper := operators[input[idx+1]]; isNextOper {
				return fmt.Errorf("%v: 2 operators in a row on the position %v\n%v",
					in, idx+1, genArrowHelp(idx+1))
			}
		}
	}
	operExists := false
	delete(operators, '.')
	for _, c := range input {
		if _, ok := operators[c]; ok {
			operExists = true
		}
	}
	if !operExists {
		return fmt.Errorf("%v: an expression does not contains any operators", input)
	}
	return
}

func main() {
	fmt.Println(inputValidate(in))
}
