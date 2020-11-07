package main

import "testing"

func BenchmarkMain(b *testing.B) {
	main()
	b.ResetTimer()
}