package decimal_test

import (
	"math"
	"testing"

	decimal "our-projectInGO/src"
)

func TestDecimalConstructor(t *testing.T) {
	// String integer
	d, err := decimal.New("123456789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Sign() != 1 {
		t.Fatalf("expected sign 1, got %d", d.Sign())
	}

	// Negative decimal
	d, err = decimal.New("-123.456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Sign() != -1 {
		t.Fatalf("expected sign -1, got %d", d.Sign())
	}

	// Negative zero float
	negZero := math.Copysign(0.0, -1.0)
	d, err = decimal.New(negZero)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Sign() != -1 {
		t.Fatalf("expected sign -1 for -0.0, got %d", d.Sign())
	}

	// Hex string
	d, err = decimal.New("0xff")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Sign() != 1 {
		t.Fatalf("expected sign 1 for 0xff, got %d", d.Sign())
	}

	// NaN
	d, err = decimal.New("NaN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Sign() != 0 {
		t.Fatalf("expected sign 0 for NaN, got %d", d.Sign())
	}

	// Infinity
	d, err = decimal.New("Infinity")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Sign() != 1 {
		t.Fatalf("expected sign 1 for Infinity, got %d", d.Sign())
	}
}

func TestIsDecimal(t *testing.T) {
	d, _ := decimal.New(42)
	if !decimal.IsDecimal(d) {
		t.Fatalf("expected IsDecimal true for Decimal object")
	}
	if decimal.IsDecimal("42") {
		t.Fatalf("expected IsDecimal false for string")
	}
}
