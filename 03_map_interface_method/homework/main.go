// 0) Легитимное ли выражение
// * Количество открывающих скобок равно количеству закрывающих.
// * Целая часть числа отделена от дробной с помощью точки
// * В строке присутствуют только допустимые символы
// 1) Токенизируем строку
// 2) Строим дерево из токенов на основе приоритетов.
//		2.1) Поиск точки «перегиба» арифметического выражения, являющийся оператором (арифметическим действием) и имеющем минимальное значение приоритета по отношению к другим операторам.
// Обходим дерево вглубину, применяя операции в каждом корневом узле на левом и правом узлах рекурсивно

package main

import (
	"fmt"
	"github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/rpn_calc"
	// "github.com/andreipimenov/golang-training-2021/03_map_interface_method/homework/bt_calc"
)

// If it is a number, terminate the recursion
/*
	bracketsFlag: Identifies the variable, used to determine whether it is in brackets
	addSubPos: record the position of the last plus and minus signs
	mulDivPos: record the position of the last multiplication sign and division sign
*/

// ИЗ учебника новосиба
//	Шаг 1. В строке найти последнюю операцию с наименьшим приоритетом.
//	Пусть ей соответствует позиция k.
//	Нужно пропускать выражения в скобках. Ввести счетчик открытых скобок

//	Шаг 2: Создать узел, содержащий знак этой операции, и рекурсивно
//	сформировать левого и правого потомка на Шаге 3.

//  Шаг 3: Левый потомок определяется подстрокой строки от k+1до N. Правый
//  потомок определяется по подстроке от 1 до k-1.
// Если длина подстроки равна 1, то надо создать новую вершину с числовым или буквенным значением

func main() {
	// calculator := bt_calc.NewCalculator()
	// fmt.Println(calculator.Calculate("1+1"))

	calculatorRPNMethod := rpn_calc.NewCalculator()
	fmt.Println(calculatorRPNMethod.Calculate("20/2-(2+2*3)"))
	fmt.Println(calculatorRPNMethod.Calculate("20/2-(2+2*3)="))
	fmt.Println(calculatorRPNMethod.Calculate("2 + (3 * 8) - (4 + (48 / (4 + 2)) * 6)"))
	fmt.Println(calculatorRPNMethod.Calculate("((81 * 6) /42+ (3-1))"))
	fmt.Println(calculatorRPNMethod.Calculate("-5/2"))

}

// 2+2*2
// (2+2)*2
// 20/2-(2+2*3)=
// ((81 * 6) /42+ (3-1))
// 2 + (3 * 8) - (4 + (48 / (4 + 2)) * 6)
/*
Additionally:

Add validation for your app to check if no symbols other than digits and allowed operators are passed as an argument to Calculate(string) method

Use structs and methods, try to implement the task using tree-based approach and stack-based approach
*/
