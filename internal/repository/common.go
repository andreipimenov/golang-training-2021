package repository

import "strings"

func SplitKey(key string) (string, string) {
	x := strings.Split(key, "_")
	if len(x) != 2 {
		return "", ""
	}
	return x[0], x[1]
}
