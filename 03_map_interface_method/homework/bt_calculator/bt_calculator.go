/*
	Пакет для парсинга математических выражений и
	вычисления результата методом построения бинаного дерева
*/

package bt_calculator

type Calculator struct {
}

func (cal *Calculator) Calculate(expression string) float64 {
	return 0.0
}

func NewCalculator() *Calculator {
	return new(Calculator)
}
