package parser

import (
	"fmt"
	"strings"
	"testing"
)

func fatal(t *testing.T, msgs ...interface{}) {
	s := "\n"
	for _, msg := range msgs {
		s += fmt.Sprint(msg)
	}
	t.Fatal(strings.TrimSuffix(s, "\n"))
}

func TestGoodExprs(t *testing.T) {
	cases := []struct {
		expr  string
		sExpr string
	}{
		{"2 + 1", "(+ 2 1)"},
		{"2 + -1", "(+ 2 (- 1))"},
		{"2 * ((4 - (2)) * 3)", "(* 2 (* (- 4 2) 3))"},
		{"45 + 5 * 4 / 6 + 7", "(+ (+ 45 (/ (* 5 4) 6)) 7)"},
		{"20/2-(2+2*3)=", "(- (/ 20 2) (+ 2 (* 2 3)))"},
		{" (4 / 7)--10.89 =", "(- (/ 4 7) (- 10.89))"},
	}

	for i, c := range cases {
		sExpr, err := Parse(c.expr)
		if err != nil {
			fatal(t,
				"i: ", i, "\n",
				err,
			)
		}
		if got, want := sExpr.String(), c.sExpr; got != want {
			fatal(t,
				"i: ", i, "\n",
				"got: ", got, "\n",
				"want: ", want, "\n",
			)
		}
	}
}

func TestBadExprs(t *testing.T) {
	cases := []struct {
		expr       string
		errMessage string
		errSpan    string
	}{
		{"45 + +", expectedNumberOrLParen, "+"},
		{"(", expectedNumberOrLParen, " "},
		{"(7 889", expectedInfixOpOrRParen, "889"},
		{" ( 7 * 8", expectedInfixOpOrRParen, " "},
		{" 7  ) ", expectedInfixOpOrExprEnd, ")"},
		{"((", expectedNumberOrLParen, " "},
		{"", expectedNumberOrLParen, " "},
		{"=", expectedNumberOrLParen, " "},
		{"2 abc", expectedInfixOpOrExprEnd, "abc"},
		{"2 + abc", expectedNumberOrLParen, "abc"},
		{"2 + &abc", expectedNumberOrLParen, "&"},
	}

	for i, c := range cases {
		_, err := Parse(c.expr)
		if err == nil {
			fatal(t,
				"i: ", i, "\n",
				"expr: '", c.expr, "'\n",
				"want message: ", c.errMessage,
				"want span: '", c.errSpan, "'\n",
			)
		}
		if err.Message != c.errMessage || err.Span.String() != c.errSpan {
			fatal(t,
				"i: ", i, "\n",
				"expr: '", c.expr, "'\n",
				"got message: ", err.Message, "\n",
				"want message: ", c.errMessage, "\n",
				"got span: '", err.Span.String(), "'\n",
				"want span: '", c.errSpan, "'\n",
			)
		}
	}
}
