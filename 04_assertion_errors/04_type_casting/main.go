package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/docker/distribution/uuid"
	uuid "github.com/satori/go.uuid"
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

type OtherDevice struct {
	ID uuid.UUID
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

	// _ = OtherDevice(d)

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
