package calc

import (
	"testing"
)

func TestCalc(t *testing.T) {
	cases := []struct {
		expr   string
		result float64
	}{
		{"20/2-(2+2*3)=", 2},
		{"8 + 3 * 2 / 4 + 3.4", 12.9},
		{"-----7", -7},
		{"(7) * 3", 21},
		{"4 + (3 + 7) * (1 + 2)", 34},
	}

	calc := Calculator{}

	for i, c := range cases {
		n, err := calc.Calculate(c.expr)

		if err != nil {
			t.Fatalf("i=%v\n%v", i, err)
		}
		if n != c.result {
			t.Fatalf("i=%v, want %v, got %v", i, n, c.result)
		}
	}
}
