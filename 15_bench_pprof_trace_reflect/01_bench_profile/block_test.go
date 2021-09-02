package block

import "testing"

func BenchmarkHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hash("alice", "john", 100, 8290)
	}
}

func BenchmarkHashBuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashBuf("alice", "john", 100, 8290)
	}
}
