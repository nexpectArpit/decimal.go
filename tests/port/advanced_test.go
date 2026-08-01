package decimal_test

import (
	"testing"

	decimal "our-projectInGO/src"
)

func TestSqrtAndCbrt(t *testing.T) {
	x, _ := decimal.New("16")
	sqrtVal := x.Sqrt()
	if sqrtVal.Sign() != 1 || sqrtVal.Exponent() != 0 {
		t.Fatalf("expected sqrt(16) e=0, got e=%d", sqrtVal.Exponent())
	}

	x27, _ := decimal.New("27")
	cbrtVal := x27.Cbrt()
	if cbrtVal.Sign() != 1 || cbrtVal.Exponent() != 0 {
		t.Fatalf("expected cbrt(27) e=0, got e=%d", cbrtVal.Exponent())
	}
}

func TestPow(t *testing.T) {
	x, _ := decimal.New("2")
	y, _ := decimal.New("10")
	powVal := x.Pow(y)

	// 2^10 = 1024
	if powVal.Sign() != 1 || powVal.Exponent() != 3 {
		t.Fatalf("expected 2^10 = 1024 (e=3), got e=%d", powVal.Exponent())
	}
}

func TestTrigAndHyperbolic(t *testing.T) {
	zero, _ := decimal.New("0")
	sin0 := zero.Sin()
	if !sin0.IsZero() {
		t.Fatalf("expected sin(0) == 0")
	}

	cos0 := zero.Cos()
	if cos0.Sign() != 1 || cos0.Exponent() != 0 {
		t.Fatalf("expected cos(0) == 1")
	}

	sinh0 := zero.Sinh()
	if !sinh0.IsZero() {
		t.Fatalf("expected sinh(0) == 0")
	}

	cosh0 := zero.Cosh()
	if cosh0.Sign() != 1 || cosh0.Exponent() != 0 {
		t.Fatalf("expected cosh(0) == 1")
	}
}
