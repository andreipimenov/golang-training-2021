package assert

import "fmt"

func ice(message string, meta []interface{}) {
	metaSuffix := ""
	if meta != nil {
		metaSuffix = fmt.Sprintf(" (meta: %v)", meta)
	}

	panic("internal calculator error: " + message + metaSuffix)
}

func Unreachable(meta ...interface{}) {
	ice("unreachable code", meta)
}

func True(value bool, failMeta ...interface{}) {
	if !value {
		ice("failed assertion", failMeta)
	}
}
