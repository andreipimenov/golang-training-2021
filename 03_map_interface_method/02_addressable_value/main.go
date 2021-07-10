package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	{
		// m := map[string]int{"first": 1, "second": 2, "last": 3}
		// p := &m["first"]
		// fmt.Println(p)
	}

	{
		// m := map[int]User{
		// 	1: {"John", 20},
		// 	2: {"Jane", 23},
		// }
		// m[1].Age = 25
	}

	{
		// s := []User{
		// 	{"John", 20},
		// 	{"Jane", 23},
		// }
		// s[1].Age = 33
		// fmt.Println(s)
	}

	{
		s := make([]int, 3, 3) // change capacity to 6 to check if pointer points to the same item
		copy(s, []int{1, 2, 3})
		p := &s[0]
		fmt.Printf("Address: %p Value: %v\n", p, *p)

		v := append(s, 4, 5, 6)
		*p = 33
		fmt.Println(s, v)
	}
}
