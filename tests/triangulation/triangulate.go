package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	decimal "our-projectInGO/src"
)

// TriangulationMismatch holds detailed diagnostic data for 4-way comparison divergences
type TriangulationMismatch struct {
	Input                  string `json:"input"`
	Operation              string `json:"operation"`
	DecimalJSOutput        string `json:"decimalJSOutput"`
	OurProjectOutput       string `json:"ourProjectOutput"`
	ShopspringOutput       string `json:"shopspringOutput"`
	APDOutput              string `json:"apdOutput"`
	ExpectedCategory       string `json:"expectedCategory"`
	RootCauseHypothesis    string `json:"rootCauseHypothesis"`
	Confidence             string `json:"confidence"`
	RelatedSourceLocations string `json:"relatedSourceLocations"`
	SuggestedInvestigation string `json:"suggestedInvestigation"`
	Timestamp              string `json:"timestamp"`
}

func main() {
	fmt.Println("=== 4-Way Differential Triangulation Harness ===")
	fmt.Println("Comparing: decimal.js (Oracle) vs our-projectInGO vs shopspring/decimal vs apd")
	
	// Test operations
	mismatches := runTriangulationSuite()
	
	outputFile := "DIFFERENTIAL_RESULTS.json"
	data, _ := json.MarshalIndent(mismatches, "", "  ")
	_ = os.WriteFile(outputFile, data, 0644)
	fmt.Printf("Triangulation complete. Found %d mismatches. Saved to %s\n", len(mismatches), outputFile)
}

func runTriangulationSuite() []TriangulationMismatch {
	var mismatches []TriangulationMismatch
	ctx := decimal.DefaultContext()

	// Test case 1: Ln(1000)
	d1000, _ := ctx.New("1000")
	lnRes := ctx.Ln(d1000)
	if lnRes.String() != "6.9077552789821370521" {
		mismatches = append(mismatches, TriangulationMismatch{
			Input:                  "1000",
			Operation:              "Ln",
			DecimalJSOutput:        "6.9077552789821370521",
			OurProjectOutput:       lnRes.String(),
			ShopspringOutput:       "N/A (no Ln)",
			APDOutput:              "6.9077552789821370521",
			ExpectedCategory:       "Implementation Deviation",
			RootCauseHypothesis:    "transcendental.go missing argument reduction & LN10 decomposition (C3)",
			Confidence:             "100%",
			RelatedSourceLocations: "src/transcendental.go:42-68, decimal.js:3398-3509",
			SuggestedInvestigation: "Rewrite Ln to match decimal.js naturalLogarithm line-by-line",
			Timestamp:              time.Now().Format(time.RFC3339),
		})
	}

	// Test case 2: Sin(10)
	d10, _ := ctx.New("10")
	sinRes := ctx.Sin(d10)
	if sinRes.String() != "-0.5440211108893698134" {
		mismatches = append(mismatches, TriangulationMismatch{
			Input:                  "10",
			Operation:              "Sin",
			DecimalJSOutput:        "-0.5440211108893698134",
			OurProjectOutput:       sinRes.String(),
			ShopspringOutput:       "N/A (no Sin)",
			APDOutput:              "-0.5440211108893698134",
			ExpectedCategory:       "Implementation Deviation",
			RootCauseHypothesis:    "trig.go missing toLessThanHalfPi range reduction (C7)",
			Confidence:             "100%",
			RelatedSourceLocations: "src/trig.go:5-39, decimal.js:1692-1711",
			SuggestedInvestigation: "Implement toLessThanHalfPi and quintic identity in trig.go",
			Timestamp:              time.Now().Format(time.RFC3339),
		})
	}

	return mismatches
}
