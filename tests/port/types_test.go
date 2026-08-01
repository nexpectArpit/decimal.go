package decimal_test

import (
	"testing"

	decimal "our-projectInGO/src"
)

func TestTypesAndConstants(t *testing.T) {
	if decimal.Base != 10000000 {
		t.Fatalf("expected Base 10000000, got %d", decimal.Base)
	}

	if decimal.LogBase != 7 {
		t.Fatalf("expected LogBase 7, got %d", decimal.LogBase)
	}

	if decimal.ExpLimit != 9000000000000000 {
		t.Fatalf("expected ExpLimit 9e15, got %d", decimal.ExpLimit)
	}

	if decimal.RoundHalfUp != 4 {
		t.Fatalf("expected RoundHalfUp 4, got %d", decimal.RoundHalfUp)
	}
}

func TestContextConfiguration(t *testing.T) {
	ctx := decimal.DefaultContext()
	if ctx.Precision != 20 {
		t.Fatalf("expected default precision 20, got %d", ctx.Precision)
	}
	if ctx.Rounding != decimal.RoundHalfUp {
		t.Fatalf("expected default rounding 4, got %d", ctx.Rounding)
	}

	ctx.Config(decimal.WithPrecision(50), decimal.WithRounding(decimal.RoundDown))
	if ctx.Precision != 50 {
		t.Fatalf("expected precision 50, got %d", ctx.Precision)
	}
	if ctx.Rounding != decimal.RoundDown {
		t.Fatalf("expected rounding RoundDown, got %d", ctx.Rounding)
	}

	cloned := ctx.Clone()
	if cloned.Precision != 50 {
		t.Fatalf("expected cloned precision 50, got %d", cloned.Precision)
	}
}
