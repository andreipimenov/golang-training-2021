package stack

type FloatStack []float64

func (s *FloatStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *FloatStack) Push(str float64) {
	*s = append(*s, str)
}

func (s *FloatStack) Pop() (float64, bool) {
	if s.IsEmpty() {
		return 0, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

type StringStack []string

func (s *StringStack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *StringStack) Push(str string) {
	*s = append(*s, str)
}

func (s *StringStack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

func (s *StringStack) Top() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		return (*s)[len(*s)-1], true
	}
}
