package lexer

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	l := NewLexer("-3 + 7 * &(4.2/-7.5)-10 ")
	type TestToken struct {
		value interface{}
		span  string
	}
	answer := []TestToken{
		{SUB, "-"}, {3, "3"}, {ADD, "+"}, {7, "7"}, {MUL, "*"}, {Unknown{}, "&"}, {LPAREN, "("},
		{4.2, "4.2"}, {DIV, "/"}, {SUB, "-"}, {7.5, "7.5"}, {RPAREN, ")"}, {SUB, "-"}, {10, "10"},
	}
	tokens := []TestToken{}

	{
		restBeforePeek := l.RestSpan()
		l.Peek()
		restAfterPeek := l.RestSpan()
		if restAfterPeek.Len() != restBeforePeek.Len() {
			t.Fatal()
		}
	}

	for {
		token := l.Peek()
		if token.IsEof() {
			break
		}
		l.Read()
		tokens = append(tokens, TestToken{token.Value, token.Span.String()})
	}
	answerStr := fmt.Sprint(answer)
	tokensStr := fmt.Sprint(tokens)
	if answerStr != tokensStr {
		t.Fatalf("\n%v - expected\n%v - actual\n", answerStr, tokensStr)
	}
}
