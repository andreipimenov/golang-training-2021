package main

import (
	"fmt"
	"unsafe"
)

func main() {
	v1 := t1{}
	fmt.Printf("Type: %T Size: %v\n", v1, unsafe.Sizeof(v1))

	v2 := t2{}
	fmt.Printf("Type: %T Size: %v\n", v2, unsafe.Sizeof(v2))

	// Struct with no fields has zero size
	v3 := struct{}{}
	fmt.Printf("Type: %T Size: %v\n", v3, unsafe.Sizeof(v3))

	// Read spec: https://golang.org/ref/spec#Size_and_alignment_guarantees

	fmt.Printf("Type: %T\n", v1)
	fmt.Printf("Type: %T Align: %v Offset: %v\n", v1.a, unsafe.Alignof(v1.a), unsafe.Offsetof(v1.a))
	fmt.Printf("Type: %T Align: %v Offset: %v\n", v1.b, unsafe.Alignof(v1.b), unsafe.Offsetof(v1.b))
	fmt.Printf("Type: %T Align: %v Offset: %v\n", v1.c, unsafe.Alignof(v1.c), unsafe.Offsetof(v1.c))

	fmt.Printf("Type: %T\n", v2)
	fmt.Printf("Type: %T Align: %v Offset: %v\n", v2.a, unsafe.Alignof(v2.a), unsafe.Offsetof(v2.a))
	fmt.Printf("Type: %T Align: %v Offset: %v\n", v2.c, unsafe.Alignof(v2.c), unsafe.Offsetof(v2.c))
	fmt.Printf("Type: %T Align: %v Offset: %v\n", v2.b, unsafe.Alignof(v2.b), unsafe.Offsetof(v2.b))
}

type t1 struct {
	a int8  // 1 byte
	b int64 // 8 bytes
	c int16 // 2 bytes
}

type t2 struct {
	a int8  // 1 byte
	c int16 // 2 bytes
	b int64 // 8 bytes
}
