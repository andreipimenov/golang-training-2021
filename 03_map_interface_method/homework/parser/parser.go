package parser

import (
	"calc/assert"
	"calc/lexer"
	"strings"
)

func getInfixPrecedence(t *lexer.Token) (p int) {
	op := t.Value.(lexer.Operator)
	switch op {
	case lexer.ADD, lexer.SUB:
		p = 1
	case lexer.MUL, lexer.DIV:
		p = 2
	default:
		assert.Unreachable(op)
	}
	return
}

func getPrefixPrecedence(op lexer.Operator) (p int) {
	if op == lexer.SUB {
		p = 3
	} else {
		assert.Unreachable(op)
	}
	return
}

func Parse(s string) (*SExpr, *ParseError) {
	s = strings.TrimSuffix(s, "=")
	l := lexer.NewLexer(s)
	expr, err := parseExpr(l, 0, true)
	if err != nil {
		return nil, err
	}
	if span := l.RestSpan(); span.Len() != 0 {
		assert.Unreachable(&span)
	}
	return expr, nil
}

// https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html
func parseExpr(l *lexer.Lexer, minPrecedence int, topLevel bool) (_ *SExpr, err *ParseError) {
	var leftExpr *SExpr
	switch leftToken := l.Read(); true {
	case leftToken.IsNumber():
		leftExpr = &SExpr{leftToken, nil}
	case leftToken.IsThisOperator(lexer.LPAREN):
		leftExpr, err = parseExpr(l, 0, false)
		if err != nil {
			return nil, err
		}
		rParen := l.Read()
		assert.True(rParen.IsThisOperator(lexer.RPAREN), &rParen.Span)
	case leftToken.IsThisOperator(lexer.SUB):
		operand, err := parseExpr(l, getPrefixPrecedence(lexer.SUB), topLevel)
		if err != nil {
			return nil, err
		}
		leftExpr = NewSExpr(leftToken, operand)
	default:
		return nil, &ParseError{expectedNumberOrLParen, leftToken.Span}
	}

	for {
		token := l.Peek()

		isInfixOp := token.IsInfixOperator()
		isRParen := token.IsThisOperator(lexer.RPAREN)
		if topLevel && !(isInfixOp || token.IsEof()) {
			return nil, &ParseError{expectedInfixOpOrExprEnd, token.Span}
		}
		if !topLevel && !(isInfixOp || isRParen) {
			return nil, &ParseError{expectedInfixOpOrRParen, token.Span}
		}
		if token.IsEof() || isRParen {
			break
		}
		// Here `token` is infix operator

		precedence := getInfixPrecedence(token)
		if precedence <= minPrecedence {
			break
		}
		l.Read()

		rightExpr, err := parseExpr(l, precedence, topLevel)
		if err != nil {
			return nil, err
		}
		leftExpr = NewSExpr(token, leftExpr, rightExpr)
	}

	return leftExpr, nil
}
