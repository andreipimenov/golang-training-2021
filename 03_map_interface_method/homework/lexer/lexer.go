package lexer

import (
	"calc/runer"
)

type Lexer struct {
	runer     runer.Runer
	nextToken *Token
}

func NewLexer(s string) *Lexer {
	runer := *runer.NewRuner(s)
	return &Lexer{runer, nil}
}

func (l *Lexer) Peek() *Token {
	if l.nextToken == nil {
		r := &l.runer
		eatSpace(r)
		if l.runer.Rest() != 0 {
			if n := readNumber(r); n != nil {
				l.nextToken = n
			} else if op := readOperator(r); op != nil {
				l.nextToken = op
			} else {
				l.nextToken = readUnknown(r)
			}
		} else {
			return &Token{Eof{}, Span{r.Str(), r.Pos(), r.Pos() + 1}}
		}
	}
	return l.nextToken
}

func (l *Lexer) Read() *Token {
	token := l.Peek()
	l.nextToken = nil
	return token
}

func (l *Lexer) RestSpan() Span {
	buffered := 0
	if l.nextToken != nil {
		buffered = l.nextToken.Span.Len()
	}

	s := l.runer.Str()
	return Span{s, l.runer.Pos() - buffered, len(s)}
}
