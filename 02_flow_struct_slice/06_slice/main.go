package main

import "fmt"

func main() {
	x := make([]byte, 0, 5)
	x = append(x, 0xFF, 0b00001111, 123, 'a')

	fmt.Println("Before add call", x, len(x), cap(x))
	add(x, 10)
	fmt.Println("After add call", x, len(x), cap(x))

	// x = x[:len(x)+1]
	// fmt.Println("After reslicing", x, len(x), cap(x))

	// add(x, 20)
	// fmt.Println("After second add call", x, len(x), cap(x))

	//x = x[:len(x)+1]

	// x = make([]byte, 1000, 1000)
	// fmt.Println("Big slice", x, len(x), cap(x))

	// s := x[0:3:3]
	// fmt.Println("Big slice reslicing with max capacity", s, len(s), cap(s))

	// s[0] = 9
	// fmt.Println(s, x)

	// s = append(s, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8)
	// fmt.Println(s, x)

	// Copying values
	// copy(x, []byte{1, 2, 3})
	// fmt.Println(x, len(x), cap(x))

	// copy(x[2:], []byte{9, 9, 9, 9, 9, 9, 9, 9, 9})
	// fmt.Println(x, len(x), cap(x))

	// copy(x[1:2], []byte{7})
	// fmt.Println(x, len(x), cap(x))
}

func add(x []byte, b byte) {
	x = append(x, b)
	fmt.Println("In add", x, len(x), cap(x))
}
