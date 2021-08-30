package user

import "testing"

func BenchmarkNewUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New("john", 25)
	}
}

func BenchmarkNewUserPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewPtr("john", 25)
	}
}
