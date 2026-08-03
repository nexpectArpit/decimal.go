package fuzz

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	decimal "our-projectInGO/src"
)

// TestDifferentialFuzz runs automated differential fuzzing against stdlib float64 reference models.
func TestDifferentialFuzz(t *testing.T) {
	rand.Seed(42)
	ctx := decimal.DefaultContext()
	iterations := 5000
	divergences := 0

	t.Logf("Running differential fuzz testing over %d random inputs...", iterations)

	for i := 0; i < iterations; i++ {
		valA := (rand.Float64() - 0.5) * 200.0
		valB := (rand.Float64() - 0.5) * 200.0

		dA, errA := ctx.New(valA)
		dB, errB := ctx.New(valB)
		if errA != nil || errB != nil {
			continue
		}

		// 1. Test Addition (dA + dB)
		addRes := ctx.Add(dA, dB)
		fAdd, err := addRes.Float64()
		if err == nil {
			expected := valA + valB
			if math.Abs(fAdd-expected) > 1e-10 {
				divergences++
			}
		}

		// 2. Test Subtraction (dA - dB)
		subRes := ctx.Sub(dA, dB)
		fSub, err := subRes.Float64()
		if err == nil {
			expected := valA - valB
			if math.Abs(fSub-expected) > 1e-10 {
				divergences++
			}
		}

		// 3. Test Multiplication (dA * dB)
		mulRes := ctx.Mul(dA, dB)
		fMul, err := mulRes.Float64()
		if err == nil {
			expected := valA * valB
			if math.Abs(fMul-expected) > 1e-10 && math.Abs(expected) > 1e-12 {
				relErr := math.Abs(fMul-expected) / math.Abs(expected)
				if relErr > 1e-10 {
					divergences++
				}
			}
		}

		// 4. Test Sqrt
		if valA > 0 {
			sqrtRes := ctx.Sqrt(dA)
			fSqrt, err := sqrtRes.Float64()
			if err == nil {
				expected := math.Sqrt(valA)
				if math.Abs(fSqrt-expected) > 1e-10 {
					divergences++
				}
			}
		}
		// 5. Test Sin
		if math.Abs(valA) > 2.0 {
			sinVal := ctx.Sin(dA)
			sinFloat, err := sinVal.Float64()
			expSin := math.Sin(valA)
			if err == nil {
				if math.Abs(sinFloat-expSin) > 1e-2 {
					divergences++
					t.Logf("[DIVERGENCE] Sin(%.6f) = %s, expected %.6f", valA, sinVal.String(), expSin)
				}
			}
		}

		// 6. Test Atan
		atanVal := ctx.Atan(dA)
		atanFloat, err := atanVal.Float64()
		expAtan := math.Atan(valA)
		if err == nil {
			if math.Abs(atanFloat-expAtan) > 1e-2 {
				divergences++
				t.Logf("[DIVERGENCE] Atan(%.6f) = %s, expected %.6f", valA, atanVal.String(), expAtan)
			}
		}

		// 7. Test Pow negative base with fractional exponent
		if valA < 0 {
			half, _ := ctx.New(0.5)
			powVal := ctx.Pow(dA, half)
			if !powVal.IsNaN() {
				divergences++
				t.Logf("[DIVERGENCE] Pow(%.6f, 0.5) = %s, expected NaN", valA, powVal.String())
			}
		}
	}

	// Explicit tests for known critical transcendental inputs from audit
	knownCases := []struct {
		op       string
		fn       func() error
	}{
		{
			op: "Sin(-70.438949)",
			fn: func() error {
				d, _ := ctx.New("-70.438949")
				res := ctx.Sin(d)
				f, err := res.Float64()
				exp := math.Sin(-70.438949)
				if err != nil || math.Abs(f-exp) > 1e-2 {
					return fmt.Errorf("Sin(-70.438949) = %s, expected %.6f", res.String(), exp)
				}
				return nil
			},
		},
		{
			op: "Atan(2.0)",
			fn: func() error {
				d, _ := ctx.New("2.0")
				res := ctx.Atan(d)
				f, err := res.Float64()
				exp := math.Atan(2.0)
				if err != nil || math.Abs(f-exp) > 1e-2 {
					return fmt.Errorf("Atan(2.0) = %s, expected %.6f", res.String(), exp)
				}
				return nil
			},
		},
		{
			op: "Atan(5.0)",
			fn: func() error {
				d, _ := ctx.New("5.0")
				res := ctx.Atan(d)
				f, err := res.Float64()
				exp := math.Atan(5.0)
				if err != nil || math.Abs(f-exp) > 1e-2 {
					return fmt.Errorf("Atan(5.0) = %s, expected %.6f", res.String(), exp)
				}
				return nil
			},
		},
		{
			op: "Ln(1000.0)",
			fn: func() error {
				d, _ := ctx.New("1000.0")
				res := ctx.Ln(d)
				f, err := res.Float64()
				exp := math.Log(1000.0)
				if err != nil || math.Abs(f-exp) > 1e-2 {
					return fmt.Errorf("Ln(1000.0) = %s, expected %.6f", res.String(), exp)
				}
				return nil
			},
		},
		{
			op: "Pow(-2.0, 0.5)",
			fn: func() error {
				d, _ := ctx.New("-2.0")
				half, _ := ctx.New("0.5")
				res := ctx.Pow(d, half)
				if !res.IsNaN() {
					return fmt.Errorf("Pow(-2.0, 0.5) = %s, expected NaN", res.String())
				}
				return nil
			},
		},
	}

	for _, kc := range knownCases {
		if err := kc.fn(); err != nil {
			divergences++
			t.Logf("[CRITICAL DIVERGENCE] %s: %v", kc.op, err)
		}
	}

	t.Logf("Differential fuzzing completed: %d operations verified, %d divergences.", iterations*7+len(knownCases), divergences)
	if divergences > 0 {
		t.Fatalf("Differential fuzzing failed with %d divergences", divergences)
	}
}

// FuzzDecimalArithmetic runs Go native fuzzing target for random byte input parsing.
func FuzzDecimalArithmetic(f *testing.F) {
	f.Add("123.456", "789.012")
	f.Add("-0.001", "1000000")
	f.Add("1e20", "-5e-10")

	f.Fuzz(func(t *testing.T, strA, strB string) {
		ctx := decimal.DefaultContext()
		dA, errA := ctx.New(strA)
		dB, errB := ctx.New(strB)

		if errA == nil && errB == nil {
			_ = ctx.Add(dA, dB)
			_ = ctx.Sub(dA, dB)
			_ = ctx.Mul(dA, dB)
		}
	})
}
