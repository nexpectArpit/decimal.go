package decimal

import (
	"testing"

	decimal "our-projectInGO/src"
)

func BenchmarkAdd(b *testing.B) {
	d1, _ := decimal.New("1234567890.123456789")
	d2, _ := decimal.New("9876543210.987654321")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d1.Add(d2)
	}
}

func BenchmarkMul(b *testing.B) {
	d1, _ := decimal.New("1234567890.123456789")
	d2, _ := decimal.New("9876543210.987654321")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d1.Mul(d2)
	}
}

func BenchmarkDiv(b *testing.B) {
	d1, _ := decimal.New("1234567890.123456789")
	d2, _ := decimal.New("9876543210.987654321")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d1.Div(d2)
	}
}

func BenchmarkLn(b *testing.B) {
	d1, _ := decimal.New("1000")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d1.Ln()
	}
}

func BenchmarkSin(b *testing.B) {
	d1, _ := decimal.New("10")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d1.Sin()
	}
}
