package main

import (
	"fmt"
)

type Device struct {
	ID string
}

func (d *Device) Print() {
	fmt.Println("Device ID is", d.ID)
}

func main() {
	d := Device{"12345"}
	d.Print()

	ds := map[string]Device{ // change Device to *Device to fix compilation error
		"first": {"666"},
	}
	_ = ds
	// ds["first"].Print()

	v := []Device{
		{"888"},
	}
	v[0].Print()
}
