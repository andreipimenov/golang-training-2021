package main

import (
	"errors"
	"fmt"
	"reflect"
)

func main() {
	typeOf()
	// kindOf()
	// valueOf()
	// elem()
	// set()
	// create()
	// fields()

	// a := &account{
	// 	123,
	// 	"admin",
	// 	"pwd",
	// 	true,
	// }
	// q, err := insertQuery(a)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(q)

	// p := Product{
	// 	ID:      "123-456-789",
	// 	Name:    "Tesla",
	// 	Price:   5000000,
	// 	InStock: false,
	// }
	// q, err = insertQuery(&p)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(q)
}

func typeOf() {
	var a int = 123

	t := reflect.TypeOf(a)
	fmt.Println(t)

	var b = struct{}{}

	t = reflect.TypeOf(b)
	fmt.Println(t)
}

func kindOf() {
	var hello = "hello"

	t := reflect.TypeOf(hello)
	k := t.Kind()
	fmt.Println(t, k)

	type user struct {
		age  int
		name string
	}

	u := user{
		age:  25,
		name: "John",
	}

	t = reflect.TypeOf(u)
	fmt.Println(t, t.Kind())

	var x = (*int)(nil)
	t = reflect.TypeOf(x)
	fmt.Println(t, t.Kind())

	var f = func(string) int { return 0 }
	t = reflect.TypeOf(f)
	fmt.Println(t, t.Kind())
}

func valueOf() {
	type mySlice []int
	var s = mySlice{1, 2, 3}

	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	fmt.Println(t, t.Kind(), v)

	var p = &[]string{"hello", "world"}
	t = reflect.TypeOf(p)
	v = reflect.ValueOf(p)
	fmt.Println(t, t.Kind(), v)
}

func elem() {
	n := []float64{1, 2, 3, 4, 5}
	v := reflect.ValueOf(&n)
	fmt.Println(v, v.Elem())

	errs := []error{
		errors.New("first error"),
		errors.New("second error"),
	}
	t := reflect.TypeOf(errs)
	fmt.Println(t, t.Elem())

	v = reflect.ValueOf(errs)
	fmt.Println(v)
}

func set() {
	type user struct {
		name string
		age  int
	}

	u := user{
		name: "Jane",
		age:  27,
	}

	newValue := reflect.ValueOf(user{
		name: "Nick",
		age:  35,
	})

	reflect.ValueOf(&u).Elem().Set(newValue)
	fmt.Println(u)
}

func create() {
	t := reflect.TypeOf(false)
	v := reflect.New(t)
	fmt.Println(v, v.Kind())

	v.Elem().Set(reflect.ValueOf(true))

	b := v.Elem().Interface()
	fmt.Printf("%T %v\n", b, b)

	f := reflect.New(reflect.TypeOf(func() {}))
	f.Elem().Set(reflect.ValueOf(func() {
		fmt.Println("runtime created function")
	}))

	f.Elem().Call(nil)

	x := f.Elem().Interface().(func())
	x()
}

func fields() {
	type Bot struct {
		ID     string `my-tag:"my-value1"`
		Active bool   `my-tag:"my-value2"`
	}

	b := Bot{
		ID:     "123-456",
		Active: true,
	}

	v := reflect.ValueOf(b)

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := v.Type().Field(i)
		fmt.Println(fieldType.Name, fieldValue.Type(), fieldValue, fieldType.Tag.Get("my-tag"))
	}
}
