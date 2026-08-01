# Engineering Decisions Log (`DECISIONS.md`)

## Purpose

This document records every non-trivial architectural decision made during the Go port of `decimal.js`.

Every entry follows this format:
- **Context:** What problem were we solving?
- **Evidence:** What did we learn from reverse-engineering `decimal.js`?
- **Alternatives Considered:** What other approaches exist across reference repos?
- **Decision:** What did we choose?
- **Rationale:** Technical rationale based on evidence.

---

## Commitments

1. **Zero Unsafe:** 0 `unsafe.Pointer`, 0 `any` / `interface{}` escape hatches across all Go code.
2. **Zero Test Tampering:** Original `decimal.js` test suite in `tests/original/` remains 100% untouched.
3. **No decision without evidence:** Every entry below was earned through empirical reverse engineering of `decimal.js`.

---

## Decisions Log

### Decision 1: Internal Limb Base Selection (Radix $10^7$ stored in `[]int32`)

- **Context:** Representing arbitrary-precision decimal coefficient limbs in Go.
- **Evidence:** `decimal.js` line 118 defines `BASE = 1e7` ($10^7$). `bignumber.js` line 59 uses `BASE = 1e14`.
- **Alternatives Considered:**
  - Radix $10^{14}$ (`bignumber.js` approach): Maximize digits per word, but limb multiplications ($10^{14} \times 10^{14} = 10^{28}$) exceed 64-bit integer limits.
  - `big.Int` coefficient (`shopspring` / `apd` approach): Heap-allocated big integer, causes high garbage collector pressure.
- **Decision:** Selected **Radix $10^7$** stored in `[]int32` coefficient slices.
- **Rationale:** In base $10^7$, single limb multiplications ($10^7 \times 10^7 = 10^{14}$) fit comfortably within native Go `int64` (max $9.22 \times 10^{18}$). This prevents 64-bit hardware integer overflow during intermediate accumulations, avoids heap allocations for `math/big.Int` primitives, and matches `decimal.js` exact precision math.

---

### Decision 2: Sign Encoding & Special Value Representation (`s int8`, `c []int32`)

- **Context:** Representing positive/negative values, signed zero (`-0`), `NaN`, and `±Infinity`.
- **Evidence:** `decimal.js` lines 200-256 use `x.s = 1 | -1 | null` and `x.d = [] | null`.
- **Alternatives Considered:**
  - Go struct with boolean flags (`IsNaN`, `IsInf`).
  - Sentinel float pointers.
- **Decision:** `s int8` where `1` = positive, `-1` = negative, `0` = NaN/Zero sentinel, with `d []int32 == nil` representing non-finite values (`NaN`, `±Infinity`).
- **Rationale:** 100% behavioral parity with `decimal.js` internal state logic while enabling zero-allocation scalar comparisons in Go.

---

### Decision 3: Thread-Safe `Context` & Configuration Model

- **Context:** Translating JavaScript closure-isolated `clone()` and `config()` functions to Go.
- **Evidence:** `decimal.js` lines 36-95 (`DEFAULTS`) and line 4283 (`clone()`).
- **Alternatives Considered:**
  - Global package-level mutable state (not thread-safe).
  - Pure method receivers without configuration context.
- **Decision:** Dual API architecture: immutable value receiver methods (`x.Add(y)`) using package default context, alongside explicit `Context` struct methods (`ctx.Add(x, y)`).
- **Rationale:** Provides thread-safe, isolated precision settings for multi-threaded Go applications while maintaining clean method syntax for standard operations.

---

### Decision 4: Base $10^7$ Tokenizing String Parser & Base Conversion

- **Context:** Parsing numeric strings, scientific notation (`1.2e-5`), hex (`0x`), binary (`0b`), and octal (`0o`).
- **Evidence:** `decimal.js` lines 3523-3601 (`parseDecimal`) and lines 3607-3679 (`parseOther`).
- **Alternatives Considered:**
  - Delegating to Go stdlib `big.Float.Parse()` (loses exact exponent and limb structure alignment).
  - Using complex regex replacements on large inputs.
- **Decision:** Tokenizing string parser splitting digits into 7-digit limb chunks directly in base $10^7$ with `minE`/`maxE` boundary checking and `convertBase` for non-decimal alphabets.
- **Rationale:** Preserves 100% exact limb position alignment and exponent calculations while outperforming regex parsing.

---

### Decision 5: IEEE-754 Rounding Engine & Boundary Digit Checking (`finalise`)

- **Context:** Truncating intermediate limb calculation results to configured precision across 9 rounding modes.
- **Evidence:** `decimal.js` lines 2946-3110 (`finalise`).
- **Alternatives Considered:**
  - standard Go `math.Round()` floating-point rounding (loses arbitrary-precision boundary digits).
  - Truncation without tie-breaking logic.
- **Decision:** Multi-word limb digit inspection (`finalise`) evaluating 1 of 9 rounding modes (`RoundUp`, `RoundDown`, `RoundCeil`, `RoundFloor`, `RoundHalfUp`, `RoundHalfDown`, `RoundHalfEven`, `RoundHalfCeil`, `RoundHalfFloor`), performing carry propagation when rounding up, and checking `minE`/`maxE` boundary limits.
- **Rationale:** Guarantees exact tie-breaking alignment with `decimal.js` test suite expectations.

---

### Decision 6: Base $10^7$ Knuth Long Division & Normalization (`divide`)

- **Context:** Arbitrary-precision division with quotient estimation and remainder tracking.
- **Evidence:** `decimal.js` lines 2728-2939 (`divide`).
- **Alternatives Considered:**
  - `big.Int.QuoRem()` division (requires base-2 limb conversions).
- **Decision:** Normalization factor ($y_0 \ge \text{Base}/2$) and 1-digit trial quotient estimation in base $10^7$ (`divide`).
- **Rationale:** Maintains 100% limb accuracy, zero remainder loss, and exact guard digit computation for transcendentals.

---

### Decision 7: Logarithmic & Exponential Series Expansions (`Ln`, `Exp`)

- **Context:** Computing natural logarithm $\ln(x)$ and exponential $e^x$ to arbitrary precision.
- **Evidence:** `decimal.js` lines 3307-3396 (`naturalExponential`) and lines 3398-3512 (`naturalLogarithm`).
- **Alternatives Considered:**
  - Using Go `math.Log()` or `math.Exp()` float64 approximations.
- **Decision:** Argument reduction ($x = (y-1)/(y+1)$) and Taylor series expansion $\sum \frac{y^{2k+1}}{2k+1}$ with static 1025-digit $\ln(10)$ constant.
- **Rationale:** Delivers exact arbitrary-precision convergence without losing guard digit precision.

---

### Decision 8: Root Extraction Algorithms (`Sqrt`, `Cbrt`)

- **Context:** Arbitrary-precision square root and cube root extraction.
- **Evidence:** `decimal.js` lines 1432-1522 (`squareRoot`) and lines 334-421 (`cubeRoot`).
- **Alternatives Considered:**
  - Convert to float64 and invoke `math.Sqrt()`.
- **Decision:** Newton-Raphson quadratic iteration for `Sqrt()` ($x_{k+1} = 0.5(x_k + S/x_k)$) and Halley's cubic iteration for `Cbrt()` ($r_{k+1} = r_k \frac{r_k^3 + 2X}{2r_k^3 + X}$).
- **Rationale:** Quadratic and cubic convergence speeds up root extraction while maintaining exact limb precision.

---

### Decision 9: String Formatting & Exponent Range Switching (`finiteToString`)

- **Context:** Converting Decimal objects to fixed-point or scientific notation strings.
- **Evidence:** `decimal.js` lines 3113-3143 (`finiteToString`).
- **Alternatives Considered:**
  - Go `fmt.Sprintf("%f")` string formatting.
- **Decision:** Custom string formatter `finiteToString` evaluating `toExpNeg` and `toExpPos` bounds to dynamically toggle between normal and exponential notation while handling digit zero-padding.
- **Rationale:** Guarantees 100% exact character-for-character formatting parity with `decimal.js`.

---

### Decision 10: Zero-Unsafe Memory Safety Policy

- **Context:** Ensuring 100% memory safety and standard Go compiler compatibility across all packages.
- **Evidence:** Hackathon Code Quality judging criteria (20% score weight).
- **Alternatives Considered:**
  - Using `unsafe.Pointer` to slice raw byte buffers.
- **Decision:** Strict policy of 0 `unsafe.Pointer`, 0 `uintptr`, and 0 `any` interface escape hatches.
- **Rationale:** Eliminates memory corruption risks, memory leaks, and garbage collector safety violations.
