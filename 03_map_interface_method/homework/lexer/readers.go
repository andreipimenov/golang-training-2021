package lexer

import (
	"calc/assert"
	"calc/runer"
	"strconv"
	"unicode"
)

func eatSpace(r *runer.Runer) {
	r.SkipWhile(unicode.IsSpace)
}

func isDot(c rune) bool {
	return c == '.'
}

func readNumber(r *runer.Runer) *Token {
	start := r.Pos()
	if r.SkipWhile(unicode.IsDigit) > 0 {
		dotPos := r.Pos()
		r.SkipIf(isDot)
		if r.SkipWhile(unicode.IsDigit) == 0 {
			r.SetPos(dotPos)
		}
	} else {
		return nil
	}
	after := r.Pos()

	n, err := strconv.ParseFloat(r.Str()[start:after], 64)
	if err != nil {
		return nil
	}
	return &Token{n, Span{r.Str(), start, after}}
}

func readOperator(r *runer.Runer) *Token {
	ops := map[rune]Operator{
		'+': ADD,
		'-': SUB,
		'*': MUL,
		'/': DIV,
		'(': LPAREN,
		')': RPAREN,
	}
	var op Operator
	start := r.Pos()
	ok := r.SkipIf(func(c rune) (known bool) {
		op, known = ops[c]
		return
	})
	if !ok {
		return nil
	}
	return &Token{op, Span{r.Str(), start, r.Pos()}}
}

func readUnknown(r *runer.Runer) *Token {
	start := r.Pos()
	if r.SkipWhile(unicode.IsLetter) == 0 {
		_, ok := r.Read()
		assert.True(ok, "must be called when not EOF")
	}
	return &Token{Unknown{}, Span{r.Str(), start, r.Pos()}}
}
