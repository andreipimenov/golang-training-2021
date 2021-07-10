package stack

import "testing"

func TestPop(t *testing.T) {
	s := NewStack()
	e := s.Pop()
	if e != nil {
		t.Errorf("stack should be empty, got %v", e)
	}
}

func TestPush(t *testing.T) {
	s := NewStack()
	s.Push(1)
	s.Push(1.0)
	s.Push("1.0")
	if s.Len != 3 {
		t.Errorf("wrong stack size: %v expected 3", s.Len)
	}
	e := s.Pop()
	if e != "1.0" {
		t.Errorf("got something wrong: %v expected '1.0'", e)
	}
}
