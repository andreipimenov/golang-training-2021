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

func main() {
	u := User{"John"}
	u.Greet()

	d := Device(u)
	d.HashName()
}
