package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

const (
	LBRACKET int32 = 40 // (
	RBRACKET int32 = 41 // )
	MULTIPL  int32 = 42 // *
	ADD      int32 = 43 // +
	SUBTR    int32 = 45 // -
	DIV      int32 = 47 // /
)

//DOT      int32 = 46 // .

func isAllowed(i int32) bool {
	return strings.ContainsRune(string(allowed()), i)
}

func allowed() []int32 {
	return append([]int32{LBRACKET, RBRACKET, MULTIPL, ADD, SUBTR, DIV}, digits()...)
}

func digits() []int32 {
	return []int32{48, 49, 50, 51, 52, 53, 54, 55, 56, 57}
}

type TokenType int

const (
	Number TokenType = iota
	Operation
	Brackets
)

type Tokens []Token

type Calc interface {
	Calculate() float64
}

func (t Tokens) Calculate() float64 {
	return t.BracketsRule().MulDivRule().AddSubtrRule().NumberRule()
}

func (t Tokens) AddSubtrRule() Tokens {
	result := Tokens{}
	for i, token := range t {
		if token.Type == Operation {
			if rune(token.Value[0]) == ADD || rune(token.Value[0]) == SUBTR {
				switch rune(token.Value[0]) {
				case ADD:
					return Tokens{Token{
						Type:  Number,
						Value: fmt.Sprintf("%f", result.Calculate()+t[i+1:].Calculate()),
					}}
				case SUBTR:
					return Tokens{Token{
						Type:  Number,
						Value: fmt.Sprintf("%f", result.Calculate()-t[i+1:].Calculate()),
					}}
				}
			} else {
				result = append(result, token)
			}
		} else {
			result = append(result, token)
		}
	}
	return result
}

func (t Tokens) MulDivRule() Tokens {
	result := Tokens{}
	for i, token := range t {
		if token.Type == Operation {
			if rune(token.Value[0]) == MULTIPL || rune(token.Value[0]) == DIV {
				switch rune(token.Value[0]) {
				case MULTIPL:
					return t.ProceedMul(result, i)
				case DIV:
					return t.ProceedDiv(result, i)
				}
			} else {
				result = append(result, token)
			}
		} else {
			result = append(result, token)
		}
	}
	return result
}

func (t Tokens) ProceedMul(result Tokens, i int) Tokens {
	before := Tokens{}
	after := Tokens{}
	if len(result) > 2 {
		before = result[:len(result)-1]
	}
	if len(t) > i+1 {
		after = t[i+2:]
	}
	return append(append(before, Token{
		Type:  Number,
		Value: fmt.Sprintf("%f", Tokens{t[i-1]}.Calculate()*Tokens{t[i+1]}.Calculate()),
	}), after...)
}

func (t Tokens) ProceedDiv(result Tokens, i int) Tokens {
	before := Tokens{}
	after := Tokens{}
	if len(result) > 2 {
		before = result[:len(result)-1]
	}
	if len(t) > i+1 {
		after = t[i+2:]
	}
	return append(append(before, Token{
		Type:  Number,
		Value: fmt.Sprintf("%f", Tokens{t[i-1]}.Calculate()/Tokens{t[i+1]}.Calculate()),
	}), after...)
}

func (t Tokens) BracketsRule() Tokens {
	result := Tokens{}
	for _, token := range t {
		if token.Type == Brackets {
			result = append(result, Token{
				Type:  Number,
				Value: fmt.Sprintf("%f", GetTokens(Lexer{src: token.Value}).Calculate()),
			})
		} else {
			result = append(result, token)
		}
	}
	return result
}

func (t Tokens) NumberRule() float64 {
	if len(t) > 0 {
		if s, err := strconv.ParseFloat(t[0].Value, 64); err == nil {
			return s
		}
	}
	return 0
}

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) WithLetterAdded(s string) Token {
	t.Value += s
	return t
}

type Lexer struct {
	src    string
	pos    int
	buffer Token
	tokens Tokens
}

func GetTokens(l Lexer) Tokens {
	if l.PosCurrentExist() {
		return GetTokens(l.TokenNext())
	}
	return l.tokens
}

func (l Lexer) PosCurrentExist() bool {
	return l.pos < len(l.src)
}
func (l Lexer) PosNextExist() bool {
	return l.pos+1 < len(l.src)
}

func (l Lexer) TokenNext() Lexer {
	if isAllowed(l.RuneCurrent()) {
		return l.WithTokenKnown()
	}
	fmt.Println("unknown char: '" + string(l.src[l.pos]) + "' at position: '" + strconv.Itoa(l.pos) + "'")
	return l.PosNext()
}

func (l Lexer) WithTokenKnown() Lexer {
	switch l.RuneCurrent() {
	case LBRACKET:
		return l.TokenBrackets()
	case MULTIPL, DIV, ADD, SUBTR:
		return l.TokenOperation()
	}
	return l.TokenNumber()
}

func (l Lexer) RuneCurrent() rune {
	return rune(l.src[l.pos])
}

func (l Lexer) RuneNext() rune {
	return rune(l.src[l.pos+1])
}

func (l Lexer) LetterCurrent() string {
	return string(l.src[l.pos])
}

func (l Lexer) TokenOperation() Lexer {
	return l.WithToken(l.TokenOperationNext()).PosNext()
}

func (l Lexer) TokenNumber() Lexer {
	return l.WithTokenBuffer(Token{
		Type: Number,
	}).TokenNumberNext()
}

func (l Lexer) TokenBrackets() Lexer {
	switch l.RuneCurrent() {
	case LBRACKET:
		return l.WithTokenBuffer(Token{
			Type:  Brackets,
			Value: "",
		}).PosNext().TokenBrackets()
	case RBRACKET:
		return l.WithToken(l.buffer).EmptyBuffer().PosNext()
	}
	return l.WithTokenBuffer(l.buffer.WithLetterAdded(l.LetterCurrent())).PosNext().TokenBrackets()
}

func (l Lexer) PosNext() Lexer {
	return l.WithPosNextNumber(1)
}

func (l Lexer) WithToken(token Token) Lexer {
	l.tokens = append(l.tokens, token)
	return l
}

func (l Lexer) WithPosNextNumber(positions int) Lexer {
	l.pos += positions
	return l
}

func (l Lexer) TokenOperationNext() Token {
	return Token{
		Type:  Operation,
		Value: l.LetterCurrent(),
	}
}

func (l Lexer) TokenNumberNext() Lexer {
	if l.PosNextExist() && unicode.IsDigit(l.RuneNext()) {
		return l.WithTokenBuffer(l.buffer.WithLetterAdded(l.LetterCurrent())).PosNext().TokenNumberNext()
	}
	return l.WithToken(l.buffer.WithLetterAdded(l.LetterCurrent())).EmptyBuffer().PosNext()
}

func (l Lexer) WithTokenBuffer(token Token) Lexer {
	l.buffer = token
	return l
}

func (l Lexer) EmptyBuffer() Lexer {
	return l.WithTokenBuffer(Token{})
}

//func main() {
//	fmt.Println(GetTokens(Lexer{src: "(2+2)*2"}).Calculate())
//}

func Test_Lazy(t *testing.T) {
	run := []struct {
		s string
		r float64
	}{
		{"(2+2)*2", 8},
		{"2+2*2", 6},
	}

	for _, test := range run {
		fmt.Println(test)
		r := GetTokens(Lexer{src: test.s}).Calculate()
		if r != test.r {
			t.Errorf("Expected %f but got %f", test.r, r)
			return
		}
	}
}
