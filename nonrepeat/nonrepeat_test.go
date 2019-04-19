package main

import "testing"

func BenchmarkLengthOfNonRepeatingSubStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if lengthOfNonRepeatingSubStr("123123") != 3 {
			b.Errorf("正确的值是:%d",3)
		}
	}
}
