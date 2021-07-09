package runer

import (
	"io"
	"strings"
)

type Runer struct {
	s      string
	reader strings.Reader
}

func NewRuner(s string) *Runer {
	return &Runer{s, *strings.NewReader(s)}
}

func (r *Runer) Str() string {
	return r.s
}

func (r *Runer) Pos() int {
	return int(r.reader.Size()) - r.reader.Len()
}

func (r *Runer) SetPos(start int) {
	r.reader.Seek(int64(start), io.SeekStart)
}

func (r *Runer) Rest() int {
	return r.reader.Len()
}

func (r *Runer) Read() (_ rune, ok bool) {
	c, _, err := r.reader.ReadRune()
	return c, err == nil
}

func (r *Runer) Peek() (_ rune, ok bool) {
	defer r.reader.UnreadRune()
	return r.Read()
}

func (r *Runer) SkipIf(skip func(rune) bool) (skipped bool) {
	c, ok := r.Peek()
	if !ok || !skip(c) {
		return false
	}
	r.Read()
	return true
}

func (r *Runer) SkipWhile(skip func(rune) bool) (skipped int) {
	for r.SkipIf(skip) {
		skipped++
	}
	return
}
