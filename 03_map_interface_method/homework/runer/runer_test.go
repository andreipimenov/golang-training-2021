package runer

import (
	"testing"
	"unicode"
)

func TestRuner(t *testing.T) {
	r := NewRuner("-55 + 7 * (4 / 7)-10")
	if !r.SkipIf(func(c rune) bool { return c == '-' }) {
		t.Fatal()
	}
	if r.SkipWhile(unicode.IsDigit) != 2 {
		t.Fatal()
	}
	if r.Str()[1:4] != "55 " {
		t.Fatal()
	}
	if r.Pos() != 3 {
		t.Fatal()
	}
	if r.Rest()+r.Pos() != len(r.Str()) {
		t.Fatal()
	}
	r.SetPos(1)
	r.Peek()
	if r.Pos() != 1 {
		t.Fatal()
	}
	got := ""
	for {
		c, ok := r.Peek()
		if !ok {
			break
		}
		r.Read()
		got += string(c)
	}
	if got != r.Str()[1:] {
		t.Fatalf("'%v' vs '%v'", got, r.Str()[1:])
	}
}
