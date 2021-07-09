package lexer

type Span struct {
	Str   string
	Start int
	After int
}

func (s *Span) String() string {
	return (s.Str + " ")[s.Start:s.After]
}

func (s *Span) Len() int {
	return s.After - s.Start
}
