package fuzz

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	decimal "our-projectInGO/src"
)

type FuzzDivergence struct {
	Op       string  `json:"op"`
	InputA   string  `json:"inputA"`
	InputB   string  `json:"inputB,omitempty"`
	GoResult string  `json:"goResult"`
	Expected string  `json:"expected"`
	Error    float64 `json:"relError"`
}

func main() {
	fmt.Println("=== Starting Phase 4 Differential Fuzz Testing ===")
	rand.Seed(42) // Deterministic seed for reproducible fuzzing

	startTime := time.Now()
	iterations := 0
	divergences := 0
	var firstDivergence *FuzzDivergence

	timeLimit := 5 * time.Second

	for time.Since(startTime) < timeLimit {
		iterations++
		aFloat := (rand.Float64() - 0.5) * 200.0

		ctx := decimal.DefaultContext()
		dA, errA := ctx.New(aFloat)
		if errA != nil {
			continue
		}

		// Check Ln: ln(a) vs math.Log(a) for a > 0
		if aFloat > 0.1 && aFloat < 100.0 {
			lnVal := ctx.Ln(dA)
			lnFloat, err := lnVal.Float64()
			expLn := math.Log(aFloat)
			if err == nil && !math.IsNaN(expLn) {
				relErr := math.Abs(lnFloat-expLn) / math.Abs(expLn)
				if relErr > 1e-3 {
					divergences++
					if firstDivergence == nil {
						firstDivergence = &FuzzDivergence{
							Op:       "Ln",
							InputA:   fmt.Sprintf("%.6f", aFloat),
							GoResult: lnVal.String(),
							Expected: fmt.Sprintf("%.6f", expLn),
							Error:    relErr,
						}
					}
				}
			}
		}

		// Check Sin: sin(a) vs math.Sin(a)
		if math.Abs(aFloat) > 2.0 {
			sinVal := ctx.Sin(dA)
			sinFloat, err := sinVal.Float64()
			expSin := math.Sin(aFloat)
			if err == nil {
				absDiff := math.Abs(sinFloat - expSin)
				if absDiff > 1e-2 {
					divergences++
					if firstDivergence == nil {
						firstDivergence = &FuzzDivergence{
							Op:       "Sin",
							InputA:   fmt.Sprintf("%.6f", aFloat),
							GoResult: sinVal.String(),
							Expected: fmt.Sprintf("%.6f", expSin),
							Error:    absDiff,
						}
					}
				}
			}
		}
	}

	duration := time.Since(startTime).Seconds()
	logText := fmt.Sprintf("# Phase 4 Differential Fuzzing Log\n# Run Duration: %.4fs\n# Total Iterations: %d\n# Total Divergences: %d\n", duration, iterations, divergences)

	if firstDivergence != nil {
		logText += fmt.Sprintf("# First Critical Divergence:\n#   Op: %s\n#   Input: %s\n#   Go Result: %s\n#   Expected: %s\n#   Rel Error: %.6f\n",
			firstDivergence.Op, firstDivergence.InputA, firstDivergence.GoResult, firstDivergence.Expected, firstDivergence.Error)
	}

	_ = os.WriteFile("fuzz/log.txt", []byte(logText), 0644)
	fmt.Printf("Fuzzing complete: %d iterations tested in %.2fs with %d divergences.\n", iterations, duration, divergences)
}
