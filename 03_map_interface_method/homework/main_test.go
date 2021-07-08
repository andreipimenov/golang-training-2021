package main

import "testing"

func TestExpr_Calculate(t *testing.T) {
	cases := []struct {
		in string
		want float64
	} {
		{"20/2-(2+2*3)",2},
		{"(235-(6*(23)))", 97},
		{"85+1-1+23+435", 543},
	}
	for _,c := range cases {
		got := Expr(c.in).Calculate()
		if got != c.want {
			t.Errorf("Calculate(%s) == %f, want %f", c.in, got, c.want)
		}
	}
}
