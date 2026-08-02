package decimal

import (
	"fmt"
	"math"
	"testing"

	decimal "our-projectInGO/src"
)

// Prove C1: parseOther hex float is broken
func TestAuditC1_HexFloatParsing(t *testing.T) {
	// 0x1.8 in decimal = 1.5, and 0x1.8p2 = 1.5 * 4 = 6.0
	d, err := decimal.New("0x1.8")
	if err != nil {
		t.Fatalf("Failed to parse 0x1.8: %v", err)
	}
	str := d.String()
	// decimal.js would return "1.5" but our code returns wrong value
	// because the divisor handling is a NOOP
	t.Logf("0x1.8 parsed as: %s (expected: 1.5)", str)
	if str != "1.5" {
		t.Errorf("AUDIT C1 CONFIRMED: 0x1.8 = %q, expected '1.5'", str)
	}
}

// Prove H8: Float64 double negation
func TestAuditH8_Float64DoubleNegation(t *testing.T) {
	d, err := decimal.New("-5.5")
	if err != nil {
		t.Fatalf("Failed to parse -5.5: %v", err)
	}
	f, err := d.Float64()
	if err != nil {
		t.Fatalf("Float64 error: %v", err)
	}
	t.Logf("-5.5.Float64() = %v (expected: -5.5)", f)
	if f != -5.5 {
		t.Errorf("AUDIT H8 CONFIRMED: Float64 of -5.5 = %v, expected -5.5", f)
	}
}

// Prove H7: Cmp returns wrong result for -Infinity
func TestAuditH7_CmpNegInfinity(t *testing.T) {
	negInf, _ := decimal.New("-Infinity")
	five, _ := decimal.New("5")
	result := negInf.Cmp(five)
	t.Logf("-Infinity.cmp(5) = %d (expected: -1)", result)
	if result != -1 {
		t.Errorf("AUDIT H7 CONFIRMED: -Infinity.cmp(5) = %d, expected -1", result)
	}
}

// Prove C3: Ln gives wrong result for numbers far from 1
func TestAuditC3_LnLargeNumber(t *testing.T) {
	// ln(1000) = 3 * ln(10) ≈ 6.907755...
	d, _ := decimal.New("1000")
	result := d.Ln()
	str := result.String()
	f, _ := result.Float64()
	expected := math.Log(1000.0)
	t.Logf("ln(1000) = %s (float64: %v, expected: ~%v)", str, f, expected)
	// Check if result is within 1% of expected
	if math.Abs(f-expected)/expected > 0.01 {
		t.Errorf("AUDIT C3 CONFIRMED: ln(1000) = %v, expected ~%v (>1%% deviation)", f, expected)
	}
}

// Prove C7: Sin gives wrong result for large angles
func TestAuditC7_SinLargeAngle(t *testing.T) {
	// sin(10) ≈ -0.5440211108893698
	d, _ := decimal.New("10")
	result := d.Sin()
	str := result.String()
	f, _ := result.Float64()
	expected := math.Sin(10.0)
	t.Logf("sin(10) = %s (float64: %v, expected: %v)", str, f, expected)
	if math.Abs(f-expected) > 0.01 {
		t.Errorf("AUDIT C7 CONFIRMED: sin(10) = %v, expected %v", f, expected)
	}
}

// Prove C10: Atan diverges for |x| > 1
func TestAuditC10_AtanLargeArg(t *testing.T) {
	// atan(2) ≈ 1.1071487177940904
	d, _ := decimal.New("2")
	result := d.Atan()
	str := result.String()
	f, _ := result.Float64()
	expected := math.Atan(2.0)
	t.Logf("atan(2) = %s (float64: %v, expected: %v)", str, f, expected)
	if math.Abs(f-expected) > 0.01 || math.IsInf(f, 0) || math.IsNaN(f) {
		t.Errorf("AUDIT C10 CONFIRMED: atan(2) = %v, expected %v", f, expected)
	}
}

// Prove C9: Sqrt with poor initial guess
func TestAuditC9_SqrtVeryLarge(t *testing.T) {
	// sqrt(1e20) = 1e10
	d, _ := decimal.New("100000000000000000000")
	result := d.Sqrt()
	str := result.String()
	t.Logf("sqrt(1e20) = %s (expected: 10000000000)", str)
	if str != "10000000000" {
		t.Logf("AUDIT C9: sqrt(1e20) = %s (may be correct but algorithm differs)", str)
	}
}

func TestAuditC6_PowNegativeBase(t *testing.T) {
	// (-2)^0.5 should be NaN in decimal.js
	negTwo, _ := decimal.New("-2")
	half, _ := decimal.New("0.5")
	result := negTwo.Pow(half)
	// decimal.js returns NaN for negative base with non-integer exponent
	if !result.IsNaN() {
		str := result.String()
		t.Errorf("AUDIT C6 CONFIRMED: (-2)^0.5 = %s, expected NaN", str)
	} else {
		t.Log("(-2)^0.5 = NaN (correct)")
	}
}

func TestAuditSummary(t *testing.T) {
	fmt.Println("\n=== PHASE 1 AUDIT DIFFERENTIAL PROOF TESTS ===")
	fmt.Println("These tests empirically verify findings from TRANSLATION_AUDIT.md")
	fmt.Println("See TRANSLATION_AUDIT.md for full analysis with decimal.js line references")
}
