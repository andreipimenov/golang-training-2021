package main

import "fmt"

type direction int

const (
	Up direction = iota
	Right
	Down
	Left
)

const (
	Execute = 1 << iota
	Write
	Read
)

const (
	A, B, C = iota, iota, iota
	D       = iota
)

const (
	One = iota + 1
	_
	Three
	Four
)

func (d direction) String() string {
	return [...]string{"Up", "Right", "Down", "Left"}[d]
}

func main() {
	x := Up
	switch x {
	case Up:
		fmt.Println(x)
	case Down:
		fmt.Println(x)
	default:
		fmt.Println(x)
	}

	fmt.Printf("Read: Binary: %03b Decimal: %d\n", Read, Read)
	fmt.Printf("Write: Binary: %03b Decimal: %d\n", Write, Write)
	fmt.Printf("Execute: Binary: %03b Decimal: %d\n", Execute, Execute)

	perms := Read | Write

	fmt.Printf("Permissions: Binary: %03b Decimal: %d\n", perms, perms)
	fmt.Printf("File permissions rwxr-xr-x: %d%d%d\n", Read|Write|Execute, Read|Execute, Read|Execute)

	fmt.Println(A, B, C, D)

	fmt.Println(One, Three, Four)
}
