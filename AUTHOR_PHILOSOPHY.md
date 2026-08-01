# Engineering Philosophy of MikeMcl (`decimal.js`, `bignumber.js`, `big.js`)

## Status: COMPLETE (Source Code Empirical Analysis)

This document synthesizes the engineering philosophy of Michael Mclaughlin (MikeMcl) across his three major arbitrary-precision JavaScript libraries:
1. `big.js` (Minimalist floating-point arithmetic library)
2. `bignumber.js` (Rational/integer arbitrary-precision arithmetic library)
3. `decimal.js` (Full-featured arbitrary-precision Decimal specification with transcendentals)

---

## 1. Architectural Evolution Across the 3 Libraries

```
   big.js (Base 10, Minimalist, ~6KB)
     │
     ▼
   bignumber.js (Base 10^14, High-Performance Rationals/Integers)
     │
     ▼
   decimal.js (Base 10^7, IEEE-754 Precision Contexts, Transcendentals)
```

### Key Differences & Design Intent:

- **`big.js`**: Designed for simple decimal math (financial calculations). Simple base 10 digits array.
- **`bignumber.js`**: Designed for speed on large integers and rational division. Uses **$10^{14}$** (`1e14`) radix to maximize digits stored per array slot. Does not include trigonometric or exponential functions.
- **`decimal.js`**: Designed for full scientific computing with arbitrary precision, supporting trigonometric (`sin`, `cos`, `tan`), logarithms (`ln`, `log`), powers (`pow`), and roots (`sqrt`, `cbrt`). Downscaled radix to **$10^7$** (`1e7`) because multi-limb series multiplications in $10^{14}$ would exceed JavaScript's 53-bit exact double float limit ($2^{53}-1 = 9 \times 10^{15}$).

---

## 2. Core Architectural Patterns Used by MikeMcl

1. **Closure-Isolated Constructor Scope (`clone()`)**:
   Constructor configuration defaults (`precision`, `rounding`, `toExpNeg`, `toExpPos`, `minE`, `maxE`) are stored in scope variables. `clone()` creates a new closure with independent configuration.

2. **Stateless Core Math Functions**:
   Public prototype methods (`.add()`, `.div()`, `.pow()`) delegate directly to pure internal helper functions (`divide()`, `naturalExponential()`, `taylorSeries()`).

3. **Explicit Guard Digits & Tail Checking**:
   Intermediate calculations compute 4 to 12 extra digits (guard digits) before calling `finalise()`. `checkRoundingDigits()` checks trailing nines (`99999`) or zeros (`00000`) near rounding boundaries to prevent precision truncation errors.

4. **Normalization & Trailing Zero Truncation**:
   Coefficients are stored normalized (no leading zero limbs). Trailing zero limbs are stripped during `finalise()`.

---

## 3. Translation Principles for `our-projectInGO`

- **Preserve Radix $10^7$**: Storing limbs as `[]int32` in base $10^7$ allows single-limb products ($10^{14}$) to fit safely within native Go `int64` ($9.22 \times 10^{18}$), keeping arithmetic fast and stack-allocated.
- **Thread-Safe Context Pattern**: Translate JS closure `clone()` into Go `Context` objects (`ctx.Add(x, y)`), preserving configuration isolation without global mutation.
- **Strict Sentinel Checks**: Represent `NaN`, `+Infinity`, `-Infinity`, and `-0` via explicit sign and coefficient slice state flags matching `decimal.js` semantics.
