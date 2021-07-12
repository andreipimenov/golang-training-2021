package bt_calc

import (
	"strconv"
)

// The task is implement the interface below
type Calc interface {
	Calculate(expression string) float64
}

func NewCalculator() *Calculator {
	return new(Calculator)
}

// Binary tree Core
// Node represents the components of binary tree search
type Node struct {
	subExpression string
	Operation     string
	NumberValue   float64
	isLeaf        bool
	LeftChild     *Node
	RightChild    *Node
}

type Calculator struct{}

func (calc *Calculator) Calculate(expression string) float64 {
	expression = trimBrackets(expression)
	root := buildExpressionTree(expression)
	return root.Calc()
}

func trimBrackets(expr string) string {
	in := []rune(expr)
	for {
		if in[0] == '(' && in[len(in)-1] == ')' {
			in = in[1 : len(in)-1]
		} else {
			break
		}
	}
	return string(in)
}

func buildExpressionTree(expression string) *Node {
	newNode := &Node{}
	expression = trimBrackets(expression)
	newNode.subExpression = expression

	parsedNumberValue, err := strconv.ParseFloat(expression, 64)
	if err == nil {
		newNode.NumberValue = parsedNumberValue
		newNode.isLeaf = true
		return newNode
	}

	splittedExpr := []byte(expression)
	infIdx := findInflectionPoint(splittedExpr)
	newNode.Operation = string(expression[infIdx])

	leftSubData := expression[:infIdx]
	if leftSubData != "" {
		newNode.LeftChild = buildExpressionTree(leftSubData)
	}
	rightSubData := expression[infIdx+1:]

	if rightSubData != "" {
		newNode.RightChild = buildExpressionTree(rightSubData)
	}
	return newNode
}

// ИЗ учебника Нижнего Новгорода
//	Шаг 1. В строке найти последнюю операцию с наименьшим приоритетом.
//	Пусть ей соответствует позиция k.
//	Нужно пропускать выражения в скобках. Ввести счетчик открытых скобок

//	Шаг 2: Создать узел, содержащий знак этой операции, и рекурсивно
//	сформировать левого и правого потомка на Шаге 3.

//  Шаг 3: Левый потомок определяется подстрокой строки от k+1до N. Правый
//  потомок определяется по подстроке от 1 до k-1.
//  Если длина подстроки равна 1, то надо создать новую вершину с числовым или буквенным значением
var operationPriority = map[string]int{
	"+": 3,
	"-": 3,
	"*": 2,
	"/": 2,
}

func findInflectionPoint(expr []byte) int {
	// Ищем такой оператор (арифметическое действие), значение которого максимально.
	// Если таковых больше одного - выбираем последний (правый) из них
	priorExpr := make([]int, len(expr))
	bracketsCounter := 0
	for nowIdx, nowValue := range expr {
		if expr[nowIdx] == '(' {
			bracketsCounter += 3
		}
		if expr[nowIdx] == ')' {
			bracketsCounter -= 3
		}
		if priority, ok := operationPriority[string(nowValue)]; ok {
			priorExpr[nowIdx] = priority + bracketsCounter
		} else {
			priorExpr[nowIdx] = 0
		}
	}

	var ans, max int
	for i := len(expr) - 1; i > 0; i-- {
		if priorExpr[i] > max {
			ans = i
			max = priorExpr[i]
		}
	}
	return ans
}

func (node *Node) Calc() float64 {
	// Максимально просто:
	// Если это лист - возвращаем число
	// Иначе считаем, вызывая Calc у детей
	if node.isLeaf {
		return node.NumberValue
	} else {
		switch node.Operation {
		case "+":
			return node.LeftChild.Calc() + node.RightChild.Calc()
		case "-":
			return node.LeftChild.Calc() - node.RightChild.Calc()
		case "*":
			return node.LeftChild.Calc() * node.RightChild.Calc()
		case "/":
			return node.LeftChild.Calc() / node.RightChild.Calc()
		}
	}
	return node.NumberValue
}

/*
Used materials:
Videos:
	"3.12 Expression trees | Binary Expression Tree | Data structures" https://www.youtube.com/watch?v=2Z6g3kNymd0&ab_channel=Jenny%27slecturesCS%2FITNET%26JRF
	"Data Structures in Golang - Binary Search Tree" - https://www.youtube.com/watch?v=-oYitelECuQ&t=4s&ab_channel=JunminLee
Articles:
	"Алгоритм парсинга арифметических выражений" - https://habr.com/ru/post/263775/
	"Дерево поиска, наивная реализация" - https://neerc.ifmo.ru/wiki/index.php?title=%D0%94%D0%B5%D1%80%D0%B5%D0%B2%D0%BE_%D0%BF%D0%BE%D0%B8%D1%81%D0%BA%D0%B0,_%D0%BD%D0%B0%D0%B8%D0%B2%D0%BD%D0%B0%D1%8F_%D1%80%D0%B5%D0%B0%D0%BB%D0%B8%D0%B7%D0%B0%D1%86%D0%B8%D1%8F
	"Программная реализация бинарных поисковых деревьев" - https://acm.bsu.by/wiki/%D0%9F%D1%80%D0%BE%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%BD%D0%B0%D1%8F_%D1%80%D0%B5%D0%B0%D0%BB%D0%B8%D0%B7%D0%B0%D1%86%D0%B8%D1%8F_%D0%B1%D0%B8%D0%BD%D0%B0%D1%80%D0%BD%D1%8B%D1%85_%D0%BF%D0%BE%D0%B8%D1%81%D0%BA%D0%BE%D0%B2%D1%8B%D1%85_%D0%B4%D0%B5%D1%80%D0%B5%D0%B2%D1%8C%D0%B5%D0%B2
	"ВВЕДЕНИЕ В СТРУКТУРЫ ДАННЫХ" by Е.А. Кумагина / Н.Н. Чернышова  - http://www.unn.ru/books/met_files/struct.pdf
*/
