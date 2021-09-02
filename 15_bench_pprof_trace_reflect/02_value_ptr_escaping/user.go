package user

type User struct {
	Name string
	Age  int
}

//go:noinline
func New(name string, age int) User {
	return User{
		Name: name,
		Age:  age,
	}
}

//go:noinline
func NewPtr(name string, age int) *User {
	return &User{
		Name: name,
		Age:  age,
	}
}
