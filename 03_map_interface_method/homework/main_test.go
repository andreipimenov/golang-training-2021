package main

import (
	"testing"
)

func TestWitoutBrackets(t *testing.T) {

	val := MyFloat(2)

	res,err:=val.Calculate("2+3/5=")
	if err!=nil || res!=2.6{
		t.Errorf("expect: %v got: %v. expression: %v. err: %v",2.6,res,"2+3/5=",err)
	}

	expression:="2+3*/5="
	res,err=val.Calculate(expression)
	if err==nil{
		t.Errorf("expression with error calculatod. expr: %v",expression)
	}
}

func TestBrackets(t *testing.T) {

	val := MyFloat(0)
	res,err:=val.Calculate("((2+3)*5+3)*2")
	if err!=nil || res!=56{
		t.Errorf("expect: %v got: %v. expression: %v. err: %v",56,res,"((2+3)*5+3)*2",err)
	}
}

func TestMoreBrackets(t *testing.T) {

	val := MyFloat(0)
	res,err:=val.Calculate("(((((-2)+(-3))*(-5))+(-3))*(-2/((-4)+(-6))))")
	if err!=nil || res!=4.4{
		t.Errorf("expect: %v got: %v. expression: %v. err: %v",56,res,"(((((-2)+(-3))*(-5))+(-3))*(-2/((-4)+(-6))))",err)
	}
}

func TestCalc(t *testing.T) {

	val := MyFloat(0)
	res,err:=val.Calculate("((-2+3)*5+3)*2")
	if err!=nil || res!=16{
		t.Errorf("expect: %v got: %v. expression: %v. err: %v",16,res,"((-2+3)*5+3)*2",err)
	}
}

func TestGarbageExpression(t *testing.T) {

	val := MyFloat(2)
	expr:="324ueou u.,3u "
	res,err:=val.Calculate(expr)
	if err==nil || res!=0{
		t.Errorf("trying to count the garbage expression. expr: %v resut: %v",res,expr)
	}
}



func TestDifferentFloatFormat(t *testing.T) {
	val := MyFloat(2)
	expr:="2.7+0.3*(5.33+7.98*(0.37+2))="
	res,err:=val.Calculate(expr)
	if err!=nil || res!=9.97278{
		t.Errorf("expect: %v got: %v. expression: %v. err: %v",9.97278,res,expr,err)
	}
}

func TestSpaceRemover(t *testing.T) {
	val := MyFloat(2)
	expr:="2   .7+0.  3*  (  5.   33+  7.9   8*(0.3  7  +2)  )="
	res,err:=val.Calculate(expr)
	if err!=nil || res!=9.97278{
		t.Errorf("expect: %v got: %v. expression: %v. err: %v",9.97278,res,expr,err)
	}
}

func TestMinusAsOperatorAndSign(t *testing.T) {
	val := MyFloat(2)
	expr:="-2.7+(-0.3)*(-5.33+7.98*(0.37-2))="
	res,err:=val.Calculate(expr)
	if err!=nil || res!=2.80122{
		t.Errorf("expect: %v got: %v. expression: %v. err: %v",2.80122,res,expr,err)
	}
}

func TestInTask(t *testing.T) {
	val := MyFloat(2)
	expr:="20/2-(2+2*3)="
	res,err:=val.Calculate(expr)
	if err!=nil || res!=2{
		t.Errorf("expect: %v got: %v. expression: %v. err: %v",2,res,expr,err)
	}
}