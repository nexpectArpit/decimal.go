# Comparative Analysis across Reference Repositories

## Status: COMPLETE (Source Code Empirical Comparison)

This document contrasts how `decimal.js` solves each mathematical problem compared to `bignumber.js`, `cockroachdb/apd`, `shopspring/decimal`, and Go standard library `math/big`.

---

## 1. Representation & Radix

| Feature | `decimal.js` | `bignumber.js` | `cockroachdb/apd` | `shopspring/decimal` | `math/big` (`big.Int`) | Go Adaptation Decision for `our-projectInGO` |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Radix / Base** | **$10^7$** (`1e7`) | $10^{14}$ (`1e14`) | Base 10 (`apd.BigInt`) | Base 10 (`big.Int` coefficient) | Base $2^{64}$ (uint64 words) | **$10^7$ (`int32` limbs)**. Matches `decimal.js` spec exactly. Single limb product ($10^{14}$) fits in Go `int64`. |
| **Sign Encoding** | `1` (pos), `-1` (neg), `null` (NaN) | `1` (pos), `-1` (neg), `null` (NaN) | `bool` (`Negative`) + Flags | `bool` (`value.Sign() < 0`) | `bool` (`neg`) | **`int8`** (`1` = positive, `-1` = negative, `0` = NaN/Zero sentinel). Zero allocation. |
| **Exponent** | Base-10 exponent `int` | Base-10 exponent `int` | `int32` Exponent | `int32` Exponent | N/A (integers) | **`int64`**. Covers $-9.22 \times 10^{18}$ to $+9.22 \times 10^{18}$, encompassing `9e15` `EXP_LIMIT`. |
| **Special Values** | `x.d = null` (NaN / Inf) | `x.c = null` (NaN / Inf) | `Form` enum (`Finite`, `Infinite`, `NaN`) | NaN not supported (panics) | N/A | **`c []int32 == nil`** with `s` sign flag. `s=0` for NaN, `s=1` for +Inf, `s=-1` for -Inf. |

---

## 2. Arithmetic & Rounding

| Operation | `decimal.js` | `bignumber.js` | `cockroachdb/apd` | `shopspring/decimal` | Adaptation Strategy |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Addition / Subtraction** | Exponent alignment + $10^7$ limb sum | Exponent alignment + $10^{14}$ limb sum | Decimal shift + `apd.BigInt` add | Decimal shift + `math/big` add | Align limb arrays by exponent difference; sum in base $10^7$; propagate carry; pass to `finalise`. |
| **Multiplication** | $O(N \cdot M)$ Schoolbook limb product in $10^7$ | $O(N \cdot M)$ Schoolbook limb product in $10^{14}$ | `apd.BigInt` multiplication | `big.Int` multiplication | Schoolbook limb product with 64-bit carry accumulator; shift exponent. |
| **Division** | Knuth-style long division with trial quotient estimation | Long division with trial quotient estimation | `Context.Quo` Newton / long division | Scaling division via `big.Int` | Implement `divide()` helper using normalization ($y_0 \ge 5 \times 10^6$) and single-digit estimation. |
| **Rounding Modes** | 9 IEEE-754 modes | 9 IEEE-754 modes | IEEE 754-2008 modes | Limited rounding modes | Map 9 modes to `RoundingMode` enum; evaluate remainder digit `rd` in `finalise()`. |

---

## 3. Advanced & Transcendental Operations

| Operation | `decimal.js` Algorithm | `cockroachdb/apd` Algorithm | `shopspring/decimal` | Adaptation Strategy |
| :--- | :--- | :--- | :--- | :--- |
| **`sqrt`** | Newton-Raphson: $x_{k+1} = 0.5(x_k + S/x_k)$ | Newton-Raphson iteration | `Sqrt()` via `big.Int` float approximation | Implement Newton-Raphson iteration matching `decimal.js` line 1432 using `divide()` and boundary nines/zeros checks. |
| **`cbrt`** | Halley's method: $r_{k+1} = r_k \frac{r_k^3 + 2X}{2r_k^3 + X}$ | N/A | N/A | Implement Halley's method matching `decimal.js` line 334. |
| **`ln` / `log`** | Argument reduction + series $\ln(\frac{1+y}{1-y})$ | Argument reduction + Taylor series | N/A | Implement argument reduction and quotient series sum matching `decimal.js` line 3398. |
| **`exp`** | Argument reduction ($x = k \ln 2 + r$) + Taylor series | Taylor series expansion | N/A | Implement Taylor series evaluation matching `decimal.js` line 3307. |
| **`sin` / `cos` / `tan`** | Argument reduction mod $2\pi$ + Taylor series | N/A | N/A | Implement Taylor series evaluator `taylorSeries()` matching `decimal.js` lines 2641 & 3687. |

---

## Summary of Architectural Verdict

Our Go port (`our-projectInGO`) will strictly follow the **$10^7$ Radix Limb representation of `decimal.js`**, avoiding heavy `big.Int` heap allocations for basic operations while maintaining 100% precision and behavioral parity across all 9 rounding modes and transcendental functions.
