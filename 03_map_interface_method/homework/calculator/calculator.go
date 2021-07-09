package calculator

import "strings"

type Calculator struct {
	input string
}

func (cal *Calculator) Calculate(expression string) float64 {
	return 0.0
}

func NewCalculator() *Calculator {
	return &Calculator{}
}

type StringsStack struct {
	data []string
	Top  string
	Len  int
}

func (s *StringsStack) Push(e string) {
	s.data = append(s.data, e)
	s.Top = s.data[len(s.data)-1]
	s.Len = len(s.data)
}

func (s *StringsStack) Pop() (string, bool) {
	if len(s.data) == 1 {
		last := s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
		s.Len = len(s.data)
		s.Top = ""
		return last, true
	}
	if len(s.data) > 0 {
		last := s.data[len(s.data)-1]
		s.data = s.data[:len(s.data)-1]
		s.Len = len(s.data)
		s.Top = s.data[len(s.data)-1]
		return last, true
	} else {
		return "", false
	}
}

func (s StringsStack) String() string {
	return strings.Join(s.data, " ")
}
