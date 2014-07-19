package models

import (
	"fmt"
	"testing"
)


func BenchmarkHello(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}
