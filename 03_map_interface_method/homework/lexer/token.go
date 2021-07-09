package lexer

type Operator int

const (
	ADD    Operator = iota // +
	SUB                    // -
	MUL                    // *
	DIV                    // /
	LPAREN                 // (
	RPAREN                 // )
)

type Number = float64

type Eof struct{}

type Unknown struct{}

type Token struct {
	Value interface{} // Operator, Number, Eof, or Unknown
	Span  Span
}

func (t *Token) IsEof() bool {
	_, ok := t.Value.(Eof)
	return ok
}

func (t *Token) IsUnknown() bool {
	_, ok := t.Value.(Unknown)
	return ok
}

func (t *Token) IsNumber() bool {
	_, ok := t.Value.(Number)
	return ok
}

func (t *Token) IsThisOperator(op Operator) bool {
	actualOp, ok := t.Value.(Operator)
	return ok && actualOp == op
}
