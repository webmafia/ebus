package ebus

import (
	"testing"
)

func Test(t *testing.T) {
	var v1 int
	var v2 int
	var v3 int64

	if typeHash(v1) != typeHash(v2) {
		t.Fatalf("%T and %T should get same hash", v1, v2)
	}

	if typeHash(v3) == typeHash(v2) {
		t.Fatalf("%T and %T should NOT get same hash", v3, v2)
	}
}

func BenchmarkTypeHash(b *testing.B) {
	for i := range b.N {
		_ = typeHash(i)
	}
}
