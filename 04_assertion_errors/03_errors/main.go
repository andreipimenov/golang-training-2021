package main

import (
	"errors"
	"fmt"
)

// Read https://blog.golang.org/go1.13-errors

var (
	NewError = fmt.Errorf("new error")
)

type ErrConnection struct{}

func (e *ErrConnection) Error() string {
	return "connection error"
}

func setupConnection() error {
	return &ErrConnection{}
}

func doSomething() error {
	return NewError
}

func wrapDoSomething() error {
	if err := doSomething(); err != nil {
		return fmt.Errorf("error doing something: %w", err)
	}
	return nil
}

func main() {
	if err := doSomething(); err != nil {
		fmt.Println(err == NewError)
		if errors.Is(err, NewError) {
			fmt.Println("this is new error")
		}
	}

	// if err := wrapDoSomething(); err != nil {
	// 	fmt.Println(err == NewError)
	// 	if errors.Is(err, NewError) {
	// 		fmt.Println("this is new error")
	// 	}
	// }

	// if err := setupConnection(); err != nil {
	// 	var connErr *ErrConnection
	// 	fmt.Println(connErr == nil)
	// 	if errors.As(err, &connErr) {
	// 		fmt.Println(connErr == nil, connErr)
	// 	}
	// }
}
