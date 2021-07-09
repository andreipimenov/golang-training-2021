package main

import (
	"bufio"
	"calc/calc"
	"calc/settings"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := []string{}
	for _, arg := range os.Args[1:] {
		if arg == "-c" || arg == "--color" {
			settings.ColoredOutput = true
		} else {
			args = append(args, arg)
		}
	}

	expr := ""
	if len(args) > 0 {
		expr = args[0]
	} else {
		reader := bufio.NewReader(os.Stdin)
		expr, _ = reader.ReadString('\n')
		expr = strings.TrimSuffix(expr, "\n")
	}

	calc := calc.Calculator{}
	n, err := calc.Calculate(expr)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Println(n)
	}
}
