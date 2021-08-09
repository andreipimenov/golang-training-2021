package mock

type Method struct {
	MethodName string
	Arguments  []interface{}
	Returns    []interface{}
}

func NewMethod() Method {
	return Method{}
}

func (m Method) WithName(n string) Method {
	m.MethodName = n
	return m
}

func (m Method) WithArguments(a ...interface{}) Method {
	m.Arguments = a
	return m
}

func (m Method) WithReturns(r ...interface{}) Method {
	m.Returns = r
	return m
}
