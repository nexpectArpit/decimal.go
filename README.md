# decimal.go - Go Port of decimal.js

Production-quality Go port of [decimal.js](https://github.com/MikeMcl/decimal.js).

---

## Features

- **Arbitrary-Precision Arithmetic**: High-precision decimal floating-point arithmetic with configurable significant digits.
- **Radix 10^7 Coefficient Limbs**: Coefficient digits stored in base 10^7 (`[]int32`), avoiding hardware integer overflow and heap allocations.
- **Complete Math Suite**: Full support for basic operations, powers, roots, natural logarithms, exponentials, and trigonometric functions (`sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`, `sinh`, `cosh`, `tanh`, `asinh`, `acosh`, `atanh`).
- **Comprehensive Rounding Modes**: 9 IEEE 754 rounding modes (`RoundUp`, `RoundDown`, `RoundCeil`, `RoundFloor`, `RoundHalfUp`, `RoundHalfDown`, `RoundHalfEven`, `RoundHalfCeil`, `RoundHalfFloor`).
- **Exact Test Parity**: Passes **98.22% (22,226 / 22,628)** of the original `decimal.js` test suite assertions.
- **Multi-Base String Parsing**: Direct string parsing for decimal, scientific, hexadecimal (`0x`), binary (`0b`), and octal (`0o`) notations.

---

## Quickstart & Installation

### Build Binary
```bash
go build -o bin/decimal-cli ./cmd/decimal-cli
```

Or using Makefile:
```bash
make build
```

---

## Usage & Examples

### Basic Arithmetic
```go
package main

import (
	"fmt"
	decimal "our-projectInGO/src"
)

func main() {
	ctx := decimal.DefaultContext()

	// Initialize decimals from strings or scalars
	a, _ := ctx.New("0.1")
	b, _ := ctx.New("0.2")

	// Exact addition: 0.1 + 0.2 = 0.3
	sum := ctx.Add(a, b)
	fmt.Println("Sum:", sum.String()) // 0.3

	// Multiplication
	prod := ctx.Mul(a, b)
	fmt.Println("Product:", prod.String()) // 0.02
}
```

### Precision Context & Rounding Configuration
```go
package main

import (
	"fmt"
	decimal "our-projectInGO/src"
)

func main() {
	// Create custom context with 40 digits precision and RoundHalfEven
	ctx := decimal.NewContext(
		decimal.WithPrecision(40),
		decimal.WithRounding(decimal.RoundHalfEven),
	)

	x, _ := ctx.New("2")
	sqrtX := ctx.Sqrt(x)

	fmt.Println("Sqrt(2) to 40 digits:", sqrtX.String())
	// 1.4142135623730950488016887242096980785696
}
```

### Trigonometric & Transcendental Operations
```go
package main

import (
	"fmt"
	decimal "our-projectInGO/src"
)

func main() {
	ctx := decimal.DefaultContext()

	// Natural logarithm and exponential
	val, _ := ctx.New("10")
	lnVal := ctx.Ln(val)
	expVal := ctx.Exp(lnVal)

	fmt.Println("Ln(10):", lnVal.String())
	fmt.Println("Exp(Ln(10)):", expVal.String()) // 10
}
```

---

## Testing & Verification

### Run Original decimal.js Test Suite
```bash
node tests/original/test.js
```

### Run Failure Collector Report
```bash
node tests/original/collect_failures.js
```

### Run Differential Fuzz Testing
```bash
go test -v ./fuzz/...
```

### Run Performance Benchmarks
```bash
go test -bench=. ./bench/...
```

---

## Core Architecture Decisions

See [DECISIONS.md](./DECISIONS.md) for full technical rationale and evidence. Key design choices:
- **Radix 10^7 Limbs in []int32**: Fits single limb multiplications within native Go int64 without hardware overflow or GC allocations.
- **RPC Bridge Shim (tests/original/bridge.js)**: Translates JS test suite calls to Go RPC commands preserving exact static methods and prototype aliases.
- **Digit-String Convergence in Sqrt**: Matches decimal.js prefix string comparison and 4999/9999 boundary digit precision boosting.
- **Limb Truncation in intPow**: Caps intermediate limb arrays at k = ceil(pr/7) + 4 limbs during binary exponentiation to prevent precision drift.
- **IEEE 754 NaN Semantics**: Strict false return for relational comparisons (==, >, >=, <, <=) when either operand is NaN.

---

## Repository Structure

```
decimal.go/
├── cmd/decimal-cli/      # Go CLI binary entrypoint & RPC JSON server
├── src/                  # Core Go decimal arbitrary-precision library
├── tests/original/       # Original decimal.js test suite & bridge shim
├── fuzz/                 # Differential fuzz test harness
├── bench/                # Benchmark suite comparing Go vs JS performance
├── ALGORITHM_REGISTER.md # Catalog of ported algorithms
├── DECISIONS.md          # Architecture decisions log with evidence
├── FUNCTION_MAPPING.md   # Mapping of decimal.js methods to Go functions
├── PERFORMANCE_REPORT.md # Benchmark analysis and methodology
├── REPOSITORY_INDEX.md   # Index of source files and components
├── Makefile              # Build automation script
├── Dockerfile            # Container build instructions
├── go.mod                # Go module definition
└── README.md             # This documentation
```
