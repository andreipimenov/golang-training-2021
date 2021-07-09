package assert

import "fmt"

const ice = "internal calculator error: "

func Unreachable(meta ...interface{}) {
	panic(fmt.Sprintf("%v%v (meta: %v)", ice, "unreachable code", meta))
}
