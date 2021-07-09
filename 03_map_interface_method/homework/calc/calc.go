package calc

import (
	"calc/assert"
	"calc/lexer"
	"calc/parser"
)

type Calc interface {
	Calculate(expr string) (float64, error)
}

type Calculator struct{}

var _ Calc = Calculator{}

func (Calculator) Calculate(expr string) (float64, error) {
	sExpr, err := parser.Parse(expr)
	if err != nil {
		return 0, err
	}
	return calculateExpr(sExpr), nil
}

func calculateExpr(expr *parser.SExpr) (n float64) {
	if len(expr.Rest) == 0 {
		return expr.Head.Value.(lexer.Number)
	}
	op := expr.Head.Value.(lexer.Operator)
	switch op {
	case lexer.ADD:
		n = calculateExpr(expr.Rest[0]) + calculateExpr(expr.Rest[1])
	case lexer.SUB:
		if len(expr.Rest) == 2 {
			n = calculateExpr(expr.Rest[0]) - calculateExpr(expr.Rest[1])
		} else {
			n = -calculateExpr(expr.Rest[0])
		}
	case lexer.MUL:
		n = calculateExpr(expr.Rest[0]) * calculateExpr(expr.Rest[1])
	case lexer.DIV:
		n = calculateExpr(expr.Rest[0]) / calculateExpr(expr.Rest[1])
	default:
		assert.Unreachable(op)
	}
	return
}
