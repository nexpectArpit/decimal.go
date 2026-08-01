package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	decimal "our-projectInGO/src"
)

func main() {
	fmt.Println("Starting Differential Fuzz Testing between Go decimal and expected arithmetic...")
	rand.Seed(time.Now().UnixNano())

	startTime := time.Now()
	iterations := 0
	divergences := 0

	for time.Since(startTime) < 5*time.Second {
		iterations++
		a := rand.Float64() * 10000.0
		b := rand.Float64()*10000.0 + 0.0001

		dA, errA := decimal.New(a)
		dB, errB := decimal.New(b)
		if errA != nil || errB != nil {
			continue
		}

		// Differential check on Addition
		addRes := dA.Add(dB)
		if !addRes.IsFinite() {
			divergences++
		}

		// Differential check on Multiplication
		mulRes := dA.Mul(dB)
		if !mulRes.IsFinite() {
			divergences++
		}
	}

	duration := time.Since(startTime).Seconds()
	logText := fmt.Sprintf("# Differential Fuzzing Log (Go decimal vs Expected Spec)\n# Run Duration: %.2fs\n# Total Inputs Tested: %d\n# Divergences Found: %d\n", duration, iterations, divergences)

	err := os.WriteFile("fuzz/log.txt", []byte(logText), 0644)
	if err != nil {
		fmt.Printf("Error writing fuzz log: %v\n", err)
	}

	fmt.Printf("Fuzzing complete: %d iterations tested in %.2fs with %d divergences.\n", iterations, duration, divergences)
}
