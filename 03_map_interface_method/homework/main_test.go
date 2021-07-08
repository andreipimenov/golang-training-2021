package main

import (
	"testing"
)

func TestInputValidate(t *testing.T) {
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
		err := inputValidate(in)
		if err == nil {
			t.Errorf("Invalid expression '%v' passed through the validation", in)
		}
	}

}
