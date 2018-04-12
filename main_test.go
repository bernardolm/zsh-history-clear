package main

import "testing"

func BenchmarkDo(b *testing.B) {
	for n := 0; n < 50; n++ {
		do()
	}
}
