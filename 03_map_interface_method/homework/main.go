package main

import (
	"fmt"
	"strconv"
	"strings"
)

type MyFloat float64

func SimpleCalc(lhs, rhs float64, op string) (float64, error) {
	switch op {
	case "*":
		return lhs * rhs, nil
	case "/":
		if rhs == 0 {
			return 0, fmt.Errorf("division by zero %v %v %v", lhs, op, rhs)
		}
		return lhs / rhs, nil
	case "+":
		return lhs + rhs, nil
	case "-":
		return lhs - rhs, nil
	default:
		return 0, fmt.Errorf("wrong operator %v", op)
	}
}

func ContainsOperator(s string) bool {
	return strings.ContainsAny(s, "+-*/()")
}

func CalcByOp(operators, expression string, opindx int) (float64, error) {

	lhs := expression[:opindx]
	lhsEdgeIndx := strings.LastIndexAny(lhs, operators)

	if lhsEdgeIndx!=-1{
		if lhs[lhsEdgeIndx]=='-' && lhsEdgeIndx!=0{
			if IsDigit(lhs[lhsEdgeIndx-1]){
				lhs = lhs[lhsEdgeIndx+1:]
			}else{
				lhs = lhs[lhsEdgeIndx:]
			}
		} else if lhs[lhsEdgeIndx]!='-'{
			lhs = lhs[lhsEdgeIndx+1:]
		}
	}

	rhs := expression[opindx+1:]
	rhsEdgeIndx := strings.IndexAny(rhs, operators)

	if rhsEdgeIndx != -1 {
		if rhsEdgeIndx == len(expression)-1 {
			return 0, fmt.Errorf("subexpression contains operator at last indx %v", expression)
		}else if rhsEdgeIndx==0 && rhs[rhsEdgeIndx]=='-'{

		} else {
			rhs = rhs[:rhsEdgeIndx]
		}
	}

	valLhs, errLhs := strconv.ParseFloat(lhs, 64)
	if errLhs != nil {
		return 0, errLhs
	}

	valRhs, errRhs := strconv.ParseFloat(rhs, 64)
	if errRhs != nil {
		return 0, errRhs
	}

	tmpres, errSmplCl := SimpleCalc(valLhs, valRhs, string(expression[opindx]))
	if errSmplCl != nil {
		return 0, errSmplCl
	}

	newexpression := strings.Replace(expression, lhs+string(expression[opindx])+rhs, fmt.Sprintf("%f", tmpres), 1)
	var nf MyFloat
	return nf.Calculate(newexpression)

}

func FindBracketIndx(s string) (int64, int64, error) {

	indxFst := strings.Index(s, "(")
	if indxFst == -1 {
		return 0, 0, fmt.Errorf("expression does not contain \"(\" %v", s)
	}
	if indxFst == len(s)-1 {
		return 0, 0, fmt.Errorf("expression does not contain \")\" %v", s)
	}

	order := 0
	for i, ch := range s[indxFst+1:] {
		if ch == '(' {
			order++
		} else if ch == ')' {
			if order == 0 {
				return int64(indxFst), int64(indxFst + i + 1), nil
			} else {
				order--
			}
		}
	}
	return 0, 0, fmt.Errorf("cant find right brackets indxs %v", s)
}

type Calc interface {
	Calculate(expression string) (float64, error)
}

func (val *MyFloat) Calculate(expression string) (float64, error) {

	expression = strings.Replace(expression, " ", "", -1)

	if strings.HasSuffix(expression, "=") {
		expression = strings.Replace(expression, "=", "", 1)
	}
	if strings.Contains(expression, "=") {
		*val = 0
		return 0, fmt.Errorf("expression contains \"=\" in body %v", expression)
	}

	if res, err := strconv.ParseFloat(expression, 64); err == nil {
		*val = MyFloat(res)
		return res, nil
	}

	if ContainsOperator(expression) {
		if strings.ContainsAny(expression, "()") {
			fstIndx, lstIndx, errBrackets := FindBracketIndx(expression)
			if errBrackets != nil {
				*val = 0
				return 0, errBrackets
			}

			subexpr := expression[fstIndx+1 : lstIndx]
			subval := MyFloat(0)
			tmpres, errCalc := subval.Calculate(subexpr)
			if errCalc != nil {
				*val = 0.0
				return 0, errCalc
			}
			newexpression := strings.Replace(expression, "("+subexpr+")", fmt.Sprintf("%f", tmpres), 1)
			result, newErr := subval.Calculate(newexpression)
			*val = MyFloat(result)
			return result, newErr

		}

		opindx := strings.IndexAny(expression, "*/")
		if opindx == 0 || opindx == len(expression)-1 {
			*val = 0.0
			return 0, fmt.Errorf("expression contains invalid operator at first or at last indx %v", expression)
		}
		if opindx != -1 {
			result, errCalcByOp := CalcByOp("*/+-", expression, opindx)
			*val = MyFloat(result)
			return result, errCalcByOp

		}

		opindx = strings.IndexAny(expression, "+-")

		if opindx == len(expression)-1 {
			*val = 0.0
			return 0, fmt.Errorf("expression contains \"+\" or \"-\" at 0 or at last indx %v", expression)
		} else if opindx == 0 {
			opindx = strings.IndexAny(expression[1:], "+-") + 1
		}


			if opindx != -1 {
				result, errCalcByOp := CalcByOp("+-", expression, opindx)
				*val = MyFloat(result)
				return result, errCalcByOp
			}
		}
	*val = 0
	return 0, fmt.Errorf("expression must contains operators. %v", expression)
}

func IsDigit(u byte) bool {
	return u >= '0' && u <= '9'
}

func main() {
	var a Calc
	val := MyFloat(2)
	a = &val
	fmt.Println(a.Calculate("((-2+3)*5+3)*2"))

}
