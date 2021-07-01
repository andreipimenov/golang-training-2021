package main

import (
	"fmt"

	"github.com/andreipimenov/golang-training-2021/01_intro/02_exported_name/user"
)

func main() {
	fmt.Println(user.Age)
	// uncomment to test if user.name is exported
	// fmt.Println(user.name)
}
