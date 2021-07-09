package parser

import (
	"calc/assert"
	"calc/lexer"
	"strings"
)

func getInfixPrecedence(op lexer.Operator) (p int) {
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
	expr, err := parseExpr(l, 0)
	if err != nil {
		return nil, err
	}
	if span := l.RestSpan(); span.Len() != 0 {
		return nil, &ParseError{unexpectedTrailingSymbols, span}
	}
	return expr, nil
}

// https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html
func parseExpr(l *lexer.Lexer, minPrecedence int) (_ *SExpr, err *ParseError) {
	var leftExpr *SExpr
	switch leftToken := l.Read(); true {
	case leftToken.IsNumber():
		leftExpr = &SExpr{leftToken, nil}
	case leftToken.IsThisOperator(lexer.LPAREN):
		leftExpr, err = parseExpr(l, 0)
		if err != nil {
			return nil, err
		}
		rParen := l.Read()
		if !rParen.IsThisOperator(lexer.RPAREN) {
			return nil, &ParseError{rParenExpected, rParen.Span}
		}
	case leftToken.IsThisOperator(lexer.SUB):
		operand, err := parseExpr(l, getPrefixPrecedence(lexer.SUB))
		if err != nil {
			return nil, err
		}
		leftExpr = NewSExpr(leftToken, operand)
	default:
		return nil, &ParseError{numberOrLParenExpected, leftToken.Span}
	}

	for {
		opToken := l.Peek()
		if opToken.IsEof() {
			break
		}
		op, ok := opToken.Value.(lexer.Operator)
		if !ok || op == lexer.LPAREN {
			return nil, &ParseError{infixOpOrExprCloserExpected, opToken.Span}
		}
		if op == lexer.RPAREN {
			break
		}

		precedence := getInfixPrecedence(op)
		if precedence <= minPrecedence {
			break
		}
		l.Read()

		rightExpr, err := parseExpr(l, precedence)
		if err != nil {
			return nil, err
		}
		leftExpr = NewSExpr(opToken, leftExpr, rightExpr)
	}

	return leftExpr, nil
}
