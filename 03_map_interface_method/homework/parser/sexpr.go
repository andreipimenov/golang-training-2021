package parser

import (
	"calc/lexer"
	"fmt"
)

// https://en.wikipedia.org/wiki/S-expression
type SExpr struct {
	Head *lexer.Token
	Rest []*SExpr
}

func NewSExpr(head *lexer.Token, rest ...*SExpr) *SExpr {
	return &SExpr{head, rest}
}

func (e *SExpr) String() string {
	s := fmt.Sprint(&e.Head.Span)
	for _, operand := range e.Rest {
		s += " " + fmt.Sprint(operand)
	}
	if e.Rest != nil {
		s = "(" + s + ")"
	}
	return s
}
