package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/calculator"
)

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
	// 3.4 В начале может быть только "-"
	operators := make(map[byte]struct{})
	for _, v := range []byte("+-*/^.") {
		operators[v] = struct{}{}
	}
	for idx, c := range input {
		if _, isOper := operators[c]; isOper {
			if idx == 0 && c != '-' {
				return fmt.Errorf("%v: not '-' operator on the 1st  position\n%v",
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
		return fmt.Errorf("%v: an expression does not contains any operators", input)
	}
	return
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

func calc(input []byte) float64 {

	// Выполняем проверку и выходим, если данные не ОК
	if err := inputValidate(input); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tokens := splitInput(input)
	rpn := reversePolishNotation(tokens)
	return calcRpn(rpn)
}

func main() {
	var in string
	e := calculator.NewCalculator()
	fmt.Println(e)

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
	// Выводим результат
	if input[len(input)-1] == '=' {
		fmt.Printf("%v%v", string(input), calc(input))
	} else {
		fmt.Printf("%v=%v", string(input), calc(input))
	}

}
