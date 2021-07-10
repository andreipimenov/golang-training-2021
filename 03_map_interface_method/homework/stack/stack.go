package stack

import (
	"fmt"
	"strings"
)

type Element struct {
	Value interface{}
	next  *Element
}

func (e Element) String() string {
	return fmt.Sprintf("%v", e.Value)
}

type Stack struct {
	Top *Element
	Len int
}

func (s *Stack) Push(value interface{}) {
	e := new(Element)
	e.Value = value
	e.next = s.Top
	s.Top = e
	s.Len++
	// s.data = append(s.data, e)
	// s.Top = s.data[len(s.data)-1]
	// s.Len = len(s.data)
}

func (s *Stack) Pop() interface{} {
	if s.Len == 0 {
		return nil
	}
	v := s.Top.Value
	s.Top = s.Top.next
	s.Len--
	return v
	// if len(s.data) == 1 {
	// 	last := s.data[len(s.data)-1]
	// 	s.data = s.data[:len(s.data)-1]
	// 	s.Len = len(s.data)
	// 	s.Top = ""
	// 	return last, true
	// }
	// if len(s.data) > 0 {
	// 	last := s.data[len(s.data)-1]
	// 	s.data = s.data[:len(s.data)-1]
	// 	s.Len = len(s.data)
	// 	s.Top = s.data[len(s.data)-1]
	// 	return last, true
	// } else {
	// 	return "", false
	// }
}

func (s *Stack) String() string {
	var vals []string
	if s.Len == 0 {
		return ""
	}
	curr := s.Top
	for i := s.Len; i > 0; i-- {
		vals = append(vals, curr.String())
		curr = curr.next
	}
	return strings.Join(vals, " ")
	// return strings.Join(s.data, " ")
}

func NewStack() *Stack {
	return new(Stack)
}
