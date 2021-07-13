package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type User struct {
	Name string
}

func (u User) Greet() {
	fmt.Println("My name is ", u.Name)
}

type Device struct {
	Name string
}

func (d Device) HashName() {
	hash := sha256.Sum256([]byte(d.Name))
	fmt.Println(hex.EncodeToString(hash[:]))
}

type HashNamer interface {
	HashName()
}

type Greeter interface {
	Greet()
}

func main() {
	u := User{"John"}
	u.Greet()

	d := Device(u)
	d.HashName()

	var i interface{} = User{"Jane"}
	x, ok := i.(Device)
	fmt.Printf("%T %v: %v\n", x, x, ok)

	z, ok := i.(Greeter)
	fmt.Println(ok)
	z.Greet()

	// z, ok := i.(HashNamer)
	// fmt.Println(ok)
	// z.HashName()
}
