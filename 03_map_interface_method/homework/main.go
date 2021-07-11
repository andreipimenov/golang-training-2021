package main

import (
	"fmt"
	"os"
	"strings"

	//calc "github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/rpn_calculator"
	calc "github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/bt_calculator"
)

// Функция выводит стрелку, которой обозначаем место неверного символа
// в тексте ошибки валидации выражения. Все.
func getArrowHelp(pos int) string {
	res := make([]rune, pos+1)
	for i := 0; i < pos; i++ {
		res[i] = ' '
	}
	res[pos] = '^'
	return string(res)
}

// Валидируем волученное выражение, возвращаем err с
// подробным описанием проблемы, если выражение некорректно
func inputValidate(input []byte) (err error) {
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
				string(input), string(c), idx, getArrowHelp(idx))
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
					string(input), idx, getArrowHelp(idx))
			}
			closing += 1
			if closing > opening {
				return fmt.Errorf("%v: closing bracket before the opening on position %v\n%v",
					string(input), idx, getArrowHelp(idx))
			}
		}
	}
	if opening > closing {
		return fmt.Errorf("%v: unclosed bracket detected on position %v\n%v",
			string(input), lastOpening, getArrowHelp(lastOpening))
	}
	// 3. Проверка операторов
	// 3.1 Проверяем, что арифметические операторы не идут подряд
	// 3.2 Проверяем, что выражение не заканчивается оператором
	// 3.3 А оператор вообще есть?
	// 3.4 В начале не может быть оператора, для отрицательных скобки плз
	operators := make(map[byte]struct{})
	for _, v := range []byte("+-*/^.") {
		operators[v] = struct{}{}
	}
	for idx, c := range input {
		if _, isOper := operators[c]; isOper {
			if idx == 0 && c == '-' {
				return fmt.Errorf("%v: please use brackets for negative operands\n%v",
					string(input), getArrowHelp(idx))
			}
			if idx == 0 {
				return fmt.Errorf("%v: operator cannot be on the 1st position\n%v",
					string(input), getArrowHelp(idx))
			}
			if idx == len(input)-1 {
				return fmt.Errorf("%v: operator on the last position\n%v",
					string(input), getArrowHelp(idx))
			}
			if _, isNextOper := operators[input[idx+1]]; isNextOper {
				return fmt.Errorf("%v: 2 operators in a row on the position %v\n%v",
					string(input), idx+1, getArrowHelp(idx+1))
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
		return fmt.Errorf("%v: an expression does not contains any operators", string(input))
	}
	return
}

func main() {
	var in string
	cal := calc.NewCalculator()

	// Если есть аргумент, считаем, что это и есть выражение
	if len(os.Args) > 1 {
		in = os.Args[1]
	} else {
		// Иначе немношк интерактивности
		fmt.Print("Enter the expression: ")
		_, err := fmt.Scanf("%v", &in)
		if err != nil {
			fmt.Printf("Your entered something wrong. Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Заменяем пробелы, убираем знак равенства в конце, если есть
	input := []byte(strings.Replace(in, " ", "", -1))
	if input[len(input)-1] == '=' {
		input = input[0 : len(input)-1]
	}
	// Валидируемся и выходим с ошибкой в случае чего
	if err := inputValidate(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	expression := string(input)
	// Выводим результат
	if input[len(input)-1] == '=' {
		fmt.Printf("%v%v", string(input), cal.Calculate(expression))
	} else {
		fmt.Printf("%v=%v", string(input), cal.Calculate(expression))
	}

}
