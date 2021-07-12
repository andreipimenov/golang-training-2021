package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func isNumber(input string) bool {
	//check if string is number
	_, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return true
	} else {
		return false
	}
}

func reformatInput(inputString string) (resultArray []string, err error) {
	/* function used for getting rid of spaces, formatting numbers consists of several digits
	in case of unexpected symbols returns error with message "incorrect input", in case of
	doubled signs returns an error*/
	var temp string
	for _, ch := range inputString {
		if isNumber(string(ch)) { //combining digits into numbers
			temp += string(ch)
		} else if string(ch) == "+" || string(ch) == "-" || string(ch) == "/" || string(ch) == "*" ||
			string(ch) == "(" || string(ch) == ")" { //if operator was found append number to slice and after that append a sign
			if len(temp) > 0 {                       //if number is not empty
				resultArray = append(resultArray, temp)
			}
			if resultArray != nil { //avoiding null pointer warning
				if resultArray[len(resultArray)-1] == string(ch) && (string(ch) != ")" && string(ch) != "(") {
					return nil, errors.New("doubled signs")
				}
			}
			resultArray = append(resultArray, string(ch))
			temp = "" //make number empty
		} else if string(ch) == "=" { //for equals sign different condition
			if len(temp) > 0 {
				resultArray = append(resultArray, temp)
			}
			temp = ""
		} else if string(ch) == " " {
		} else {
			return nil, errors.New("incorrect input")
		}
	}
	return
}

func isSign(signs map[string]int, elem string) bool {
	//check if symbol is sign, signs are kept in the map
	_, ok := signs[elem]
	return ok
}

func performOperation(integers []float64, sign string) {
	switch sign {
	case "+":
		integers[len(integers)-2] = integers[len(integers)-2] + integers[len(integers)-1]
	case "-":
		integers[len(integers)-2] = integers[len(integers)-2] - integers[len(integers)-1]
	case "/":
		integers[len(integers)-2] = integers[len(integers)-2] / integers[len(integers)-1]
	case "*":
		integers[len(integers)-2] = integers[len(integers)-2] * integers[len(integers)-1]
	}
}

type Calc interface {
	Calculate(userInput string) (float64, error)
}

type Input struct {
}

func (Input) Calculate(inputString string) (float64, error) {
	operators := map[string]int{
		"-": 0,
		"+": 0,
		"*": 1,
		"/": 1,
		"(": 2,
		")": 2,
	}
	//тестовые выражения
	//inputString = "223*231/123*(1/32*(22/2*(11+11*(23*231/123*(524)))))="
	//inputString = "23 * 123 * (15 / 21 / (123 * 54)*(522/(22*32*(12/4))))="
	var resultArray []string
	if !strings.Contains(inputString, "=") {
		return 0, errors.New("input should have '=' at the end of the statement")
	}

	formattedInput, err := reformatInput(inputString)
	if err != nil {
		return 0, err
	}

	var operatorStack stack //declare a stack for operators

	//fmt.Println(formattedInput) for checking how input formatter worked
	for _, elem := range formattedInput {
		if isSign(operators, elem) {
			if operatorStack.size() == 0 {
				operatorStack.push(elem)
			} else if elem == "(" { //if bracket is open ignoring priority weight and push signs
				operatorStack.push(elem)
			} else if elem == ")" { //if bracket is closed popping all signs from stack between these brackets and putting it to result array
				for operatorStack.peek() != "(" {
					resultArray = append(resultArray, operatorStack.pop())
				}
				operatorStack.pop()
			} else if operators[operatorStack.peek()] < operators[elem] || operatorStack.peek() == "(" { //comparing operators priorities
				operatorStack.push(elem)
			} else if operators[operatorStack.peek()] >= operators[elem] {
				resultArray = append(resultArray, operatorStack.pop())
				operatorStack.push(elem)
			}
		} else { //if string is not operator, means it is number, so adding it to result array
			resultArray = append(resultArray, elem)
		}
	}
	//if any operators were left in stack put them to result array
	for operatorStack.size() != 0 {
		resultArray = append(resultArray, operatorStack.pop())
	}
	//resultArray is a slice for reverse polish notation result
	var floats []float64 //we'll put numbers in floats, and if sign is met perform operations
	for _, elem := range resultArray {
		if isNumber(elem) {
			elemFloat, _ := strconv.ParseFloat(elem, 64)
			floats = append(floats, elemFloat)
		} else {
			performOperation(floats, elem)
			floats = floats[:len(floats)-1]
		}
	}

	if floats != nil { //null pointer check
		return floats[0], nil
	} else {
		return 0, errors.New("something went wrong with calculations")
	}
}

func main() {
	var calculator Calc = Input{}
	userInput := "23 * 123 * (15 / 21 / (123 * 54)*(522/(22*32*(12/4))))="
	res, err := calculator.Calculate(userInput)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(userInput, res)
	}
}
