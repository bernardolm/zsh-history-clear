package main

import "testing"

func benchmarkDo(i int, b *testing.B) {
	for n := 0; n < 100; n++ {
		do()
	}
}

func BenchmarkDo1(b *testing.B)           { benchmarkDo(1, b) }
func BenchmarkDo100(b *testing.B)         { benchmarkDo(100, b) }
func BenchmarkDo1000(b *testing.B)        { benchmarkDo(1000, b) }
func BenchmarkDo10000(b *testing.B)       { benchmarkDo(10000, b) }
func BenchmarkDo100000(b *testing.B)      { benchmarkDo(100000, b) }
func BenchmarkDo1000000(b *testing.B)     { benchmarkDo(1000000, b) }
func BenchmarkDo10000000(b *testing.B)    { benchmarkDo(10000000, b) }
func BenchmarkDo100000000(b *testing.B)   { benchmarkDo(100000000, b) }
func BenchmarkDo1000000000(b *testing.B)  { benchmarkDo(1000000000, b) }
func BenchmarkDo10000000000(b *testing.B) { benchmarkDo(10000000000, b) }
