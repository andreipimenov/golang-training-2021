package main

import (
	"testing"
)

func TestInvalidExpressions(t *testing.T) {
	invalidInputs := []string{
		"123.0",
		"2++",
		"20/2-(2+2*3)^2==",
		"20./2-(2+2*3)^2",
		"20/2-((2+2*3)^2",
		"20/2-()2+2*3)^2",
		"20/2-(()2+2*3)^2",
		"20/2-(2++2*3)^2=",
		"20/2-(2+2*3)^2.",
		"20/2-(2+2*3)^-2",
		"20/2-(2+2*3)^.2",
		"*20.2",
	}
	for _, in := range invalidInputs {
		err := inputValidate([]byte(in))
		if err == nil {
			t.Errorf("Invalid expression '%v' passed through the validation", in)
		}
	}
}

func TestValidExpressions(t *testing.T) {
	validInputs := []string{
		"2+2*2",
		"20/2-(2+2*3)^2",
		"20.1/2-(2+2*3)^2",
	}
	for _, in := range validInputs {
		err := inputValidate([]byte(in))
		if err != nil {
			t.Errorf("Valid expression '%v' didn't passed through the validation\n%v", in, err)
		}
	}
}
