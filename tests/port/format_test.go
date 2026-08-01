package decimal_test

import (
	"math"
	"testing"

	decimal "our-projectInGO/src"
)

func TestFormatting(t *testing.T) {
	x, _ := decimal.New("123.456")
	if x.String() != "123.456" {
		t.Fatalf("expected '123.456', got '%s'", x.String())
	}

	if x.ToFixed(2) != "123.46" {
		t.Fatalf("expected ToFixed(2) '123.46', got '%s'", x.ToFixed(2))
	}

	if x.ToExponential(2) != "1.23e+2" {
		t.Fatalf("expected ToExponential(2) '1.23e+2', got '%s'", x.ToExponential(2))
	}

	if x.ToPrecision(4) != "123.5" {
		t.Fatalf("expected ToPrecision(4) '123.5', got '%s'", x.ToPrecision(4))
	}
}

func TestFloat64Conversion(t *testing.T) {
	x, _ := decimal.New("3.141592653589793")
	fVal, err := x.Float64()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if math.Abs(fVal-3.141592653589793) > 1e-15 {
		t.Fatalf("expected ~3.141592653589793, got %v", fVal)
	}
}
