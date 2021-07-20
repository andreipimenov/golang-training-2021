package main

import "fmt"

type Worker interface {
	Do()
}

type User struct{}

func (u *User) Do() {
	fmt.Println("Doing...")
}

func main() {
	var w Worker
	fmt.Println(w == nil)

	var u *User
	fmt.Println(u == nil)
	// u.Do()

	w = u
	fmt.Println(w == nil)

	// fmt.Printf("%T %v\n", w, w)
}
