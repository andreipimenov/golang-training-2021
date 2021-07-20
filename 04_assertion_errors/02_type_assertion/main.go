package main

import "fmt"

type User struct{}

func (u *User) Greet() {
	fmt.Println("Hello!")
}

func (u *User) Other() {
	fmt.Println("Other method!")
}

func main() {
	var u interface{} = &User{}

	if _, ok := u.(interface{ Greet() }); ok {
		fmt.Println("u has Greet() method")
		// fmt.Printf("%T %v\n", x, x)
		// x.Greet()
	}

}
