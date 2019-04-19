package main

import "testing"

func TestTringle(t *testing.T) {
	tests := []struct{ a, b, c int }{
		{3, 4, 5},
		{30000, 40000, 50000},
		{5, 12, 13},
		{8, 15, 17},
		{12, 35, 37},
	}
	for _, tt := range tests {
		if calcTringle(tt.a, tt.b) != tt.c {
			t.Errorf("正确的值是:%d", tt.c)
		}
	}
}
func BenchmarkTringle(b *testing.B) {
	a := 30000
	d := 40000
	c := 50000
	for i := 0; i < b.N; i++ {
		if calcTringle(a, d) != c {
			b.Errorf("正确的值是:%d", c)
		}
	}
}
