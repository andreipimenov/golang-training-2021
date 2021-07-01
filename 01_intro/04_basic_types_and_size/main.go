package main

import (
	"fmt"
	"math/cmplx"
	"unsafe"
)

var (
	a bool       = false
	b uint64     = 1<<64 - 1
	c complex128 = cmplx.Sqrt(-5 + 12i)
	d string     = "hello"
	v rune       = 'h'
)

func main() {
	fmt.Printf("Type: %T Value: %v Size: %v\n", a, a, unsafe.Sizeof(a))
	fmt.Printf("Type: %T Value: %v Size: %v\n", b, b, unsafe.Sizeof(b))
	fmt.Printf("Type: %T Value: %v Size: %v\n", c, c, unsafe.Sizeof(c))
	fmt.Printf("Type: %T Value: %v Size: %v\n", d, d, unsafe.Sizeof(d))
	fmt.Printf("Type: %T Value: %v Size: %v\n", v, v, unsafe.Sizeof(v))
}
