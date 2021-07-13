package rpn_calculator

import (
	"strings"
	"testing"
)

func TestSplitInput(t *testing.T) {
	testData := map[string]string{
		"2+2*2":            "2 + 2 * 2",
		"20/2-(2+2*3)^2":   "20 / 2 - ( 2 + 2 * 3 ) ^ 2",
		"20.1/2-(2+2*3)^2": "20.1 / 2 - ( 2 + 2 * 3 ) ^ 2",
		"(-2)^2":           "( - 2 ) ^ 2",
	}
	for k, v := range testData {
		arr := splitInput([]byte(k))
		if strings.Join(arr, " ") != v {
			t.Errorf("Wrong splitting: got: '%v', expected: '%v'", arr, v)
		}
	}
}

func TestRPN(t *testing.T) {
	// За истину берем данные с калькулятора
	// https://abakbot.ru/online-13/181-obratnaya-polskaya-notatsiya-konvertatsiya-onlajn
	testData := map[string]string{
		"2 + 2 * 2":                    "2 2 2 * +",
		"14 / 88":                      "14 88 /",
		"20 / 2 - ( 2 + 2 * 3 ) ^ 2":   "20 2 / 2 2 3 * + 2 ^ -",
		"20.1 / 2 - ( 2 + 2 * 3 ) ^ 2": "20.1 2 / 2 2 3 * + 2 ^ -",
		"( - 2 ) ^ 2":                  "0 2 - 2 ^",
		"100.5 + 13 / 37":              "100.5 13 37 / +",
		"( 2 * 2 * 2 ) ^ 3":            "2 2 * 2 * 3 ^",
		"20 / 2 - ( 2 + 2 * 3 )":       "20 2 / 2 2 3 * + -",
	}
	for k, v := range testData {
		arr := strings.Split(k, " ")
		res := strings.Join(reversePolishNotation(arr), " ")
		if res != v {
			t.Errorf("Wrong RPN convertation: got: '%v', expected: '%v'", arr, v)
		}
	}
}

func TestCalcRPN(t *testing.T) {
	cal := NewCalculator()
	testData := map[string]float64{
		"2+2*2=":                       6.0,
		"(-2)^2=":                      4.0,
		"20/2-(2+2*3)^2=":              -54,
		"2 + 2 * 2":                    6,
		"20.1 / 2 - ( 2 + 2 * 3 ) ^ 2": -53.95,
		"( - 2 ) ^ 2":                  4,
		"20/2-(2+2*3)":                 2,
		"14/88=":                       0.1590909090909091,
		"2^(2+2)":                      16,
		"2^(-2)":                       0.25,
		"2":                            2,
		"(((2+2)))":                    4,
	}
	for k, v := range testData {
		input := []byte(strings.Replace(k, " ", "", -1))
		if input[len(input)-1] == '=' {
			input = input[0 : len(input)-1]
		}
		res := cal.Calculate(string(input))
		if res != v {
			t.Errorf("Wrong calc result of %v: got: '%v', expected: '%v'", k, res, v)
		}
	}
}
