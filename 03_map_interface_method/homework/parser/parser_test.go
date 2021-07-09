package parser

import (
	"testing"
)

func TestGoodExprs(t *testing.T) {
	cases := []struct {
		expr  string
		sExpr string
	}{
		{"45 + 5 * 4 / 6 + 7", "(+ (+ 45 (/ (* 5 4) 6)) 7)"},
		{"20/2-(2+2*3)=", "(- (/ 20 2) (+ 2 (* 2 3)))"},
		{" (4 / 7)--10.89 =", "(- (/ 4 7) (- 10.89))"},
	}

	for i, c := range cases {
		sExpr, err := Parse(c.expr)
		if err != nil {
			t.Fatalf("i=%v\n%v", i, err)
		}
		if got, want := sExpr.String(), c.sExpr; got != want {
			t.Fatalf("i=%v\n'%v'\n'%v'", i, got, want)
		}
	}
}

func TestBadExprs(t *testing.T) {
	cases := []struct {
		expr       string
		errMessage string
		errSpan    string
	}{
		{"45 + +", numberOrLParenExpected, "+"},
		{"(", numberOrLParenExpected, " "},
		{"(7 889", infixOpOrExprCloserExpected, "889"},
		{" ( 7 * 8", rParenExpected, " "},
		{" 7  ) ", unexpectedTrailingSymbols, ") "},
		{"", numberOrLParenExpected, " "},
		{"=", numberOrLParenExpected, " "},
	}

	for i, c := range cases {
		_, err := Parse(c.expr)
		if err == nil {
			t.Fatalf("i=%v '%v'", i, c.expr)
		}
		if got, want := err.Message, c.errMessage; got != want {
			t.Fatalf("i=%v\n'%v'\n'%v'", i, got, want)
		}
		if got, want := err.Span.String(), c.errSpan; got != want {
			t.Fatalf("i=%v\n'%v'\n'%v'", i, got, want)
		}
	}
}
