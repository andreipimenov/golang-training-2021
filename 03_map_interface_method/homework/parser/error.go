package parser

import (
	"calc/lexer"
	"calc/settings"
	"strings"
)

const (
	unexpectedTrailingSymbols   = "unexpected trailing symbols"
	rParenExpected              = "')' is expected"
	numberOrLParenExpected      = "number (possibly with unary '-') or '(' is expected"
	infixOpOrExprCloserExpected = "infix operator or ')' (or expression end on the top level) is expected"
)

type ParseError struct {
	Message string
	Span    lexer.Span
}

const (
	RESET = "\033[0m"
	RED   = RESET + "\033[1;31m"
	BOLD  = RESET + "\033[1m"
)

func (e *ParseError) Error() (s string) {
	reset := ""
	red := ""
	bold := ""
	if settings.ColoredOutput {
		reset = RESET
		red = RED
		bold = BOLD
	}

	s += red + "Error" + bold + ": " + e.Message + reset + "\n"
	s += e.Span.Str + "\n"
	s += strings.Repeat(" ", e.Span.Start) +
		red + strings.Repeat("^", e.Span.Len()) + reset + "\n"
	return s
}
