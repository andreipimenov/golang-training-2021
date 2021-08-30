package main

import "fmt"

type Service interface {
	GetPrice(string) int
}

func main() {
	var s Service
	fmt.Println(s.GetPrice("AAPL"))
}
