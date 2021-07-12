package main

import "fmt"

type stack struct {
	stackArr []string
}

func (s *stack) pop() (poppedElem string) {
	if len(s.stackArr) > 0 {
		poppedElem = s.stackArr[len(s.stackArr)-1]
		s.stackArr = append(s.stackArr[:len(s.stackArr)-1])
	} else {
		fmt.Println("Stack is empty")
	}
	return
}

func (s *stack) push(elem string) {
	s.stackArr = append(s.stackArr, elem)
}

func (s *stack) peek() (lastElem string) {
	return s.stackArr[len(s.stackArr)-1]
}

func (s *stack) size() int {
	return len(s.stackArr)
}
