package main

import "fmt"

type Printer interface {
	Print()
}

type User struct {
	Name string
}

func (u *User) Print() {
	fmt.Println("User name is", u.Name)
}

type Device struct {
	ID string
}

func (d Device) Print() {
	fmt.Println("Device ID is", d.ID)
}

func main() {
	ps := []Printer{
		&User{"John"}, // Impossible to pass User but *User because interface value is not addressable
	}
	ps = append(ps, Device{"12345"})

	for _, p := range ps {
		p.Print()
	}

	ps = append(ps, &Device{"0001"})
	ps[2].Print()
}
