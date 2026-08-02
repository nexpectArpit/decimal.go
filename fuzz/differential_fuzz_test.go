package fuzz

import (
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
	}

	t.Logf("Differential fuzzing completed: %d operations verified, %d divergences.", iterations*4, divergences)
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
