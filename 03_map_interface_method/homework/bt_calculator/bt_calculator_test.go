package bt_calculator

import (
	"strings"
	"testing"
)

func TestFindInflectionPoint(t *testing.T) {
	testData := map[string]int{
		"2+2*2":              1,
		"(1+10.2)^2+5*0.5-2": 16,
		"2":                  0,
	}
	for k, v := range testData {
		r := findInflectionPoint([]byte(k))
		if r != v {
			t.Errorf("wrong root in %v got: %v expected: %v", k, r, v)
		}
	}
}

func TestTrimBrackets(t *testing.T) {
	testData := map[string]string{
		"(2+2)*2":  "(2+2)*2",
		"(1+10.2)": "1+10.2",
	}
	for k, v := range testData {
		r := trimBrackets(k)
		if r != v {
			t.Errorf("wrong trimBrackets %v got: %v expected: %v", k, r, v)
		}
	}
}

func TestCalculate(t *testing.T) {
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
	cal := NewCalculator()
	for k, v := range testData {
		input := []byte(strings.Replace(k, " ", "", -1))
		if input[len(input)-1] == '=' {
			input = input[0 : len(input)-1]
		}
		res := cal.Calculate(string(input))
		if res != v {
			t.Errorf("Wrong calc result of %v: got: '%v', expected: '%v'", string(input), res, v)
		}
	}
}
