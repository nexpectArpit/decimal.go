# Algorithm Register

## Status: COMPLETE (Derived from `decimal.js` Source Code Inspection)

This register documents every algorithm used across `decimal.js`, verified with line numbers, complexity analysis, and Go implementation strategies.

---

## Arithmetic Operations

| Function | Algorithm | Evidence (`decimal.js` line) | Complexity | Dependencies | Go Plan |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `plus` / `add` | Base $10^7$ alignment, limb addition with carry propagation | Lines 1262-1358 | $O(N)$ | `finalise`, `getBase10Exponent` | Align limb arrays by exponent difference; sum array elements in base $10^7$; propagate carry; pass to `finalise`. |
| `minus` / `sub` | Sign-flipped addition or limb subtraction with borrow propagation | Lines 1024-1110 | $O(N)$ | `finalise`, `getBase10Exponent` | Compare magnitudes; subtract smaller from larger in base $10^7$; propagate borrow; pass to `finalise`. |
| `times` / `mul` | $O(N \cdot M)$ Schoolbook limb multiplication with carry accumulator | Lines 1577-1678 | $O(N \cdot M)$ | `finalise`, `getBase10Exponent` | Loop nested limb indices $i, j$; accumulate $d_1[i] \times d_2[j]$ into 64-bit int; divide by $10^7$ for carry; pass to `finalise`. |
| `dividedBy` / `div` | Base $10^7$ Knuth-style long division with quotient estimation & normalization | Lines 2728-2939 | $O(N \cdot M)$ | `multiplyInteger`, `compare`, `subtract`, `finalise` | Normalize divisor ($y[0] \ge 5 \times 10^6$); estimate trial quotient digit $k = \text{rem} / y_0$; subtract trial product; adjust remainder; scale exponent. |
| `mod` / `modulo` | Division quotient truncation: $r = a - n \times q$, where $q = \text{trunc}(a / n)$ | Lines 1133-1185 | $O(N \cdot M)$ | `divide`, `finalise` | Compute $q = \text{div}(x, y)$ with rounding mode; return $x - (y \times q)$. |
| `divToInt` | Truncated division: $q = \text{floor/trunc}(a / n)$ | Lines 478-482 | $O(N \cdot M)$ | `divide`, `finalise` | Invoke `divide` with truncated integer precision; pass to `finalise`. |

---

## Comparison & Predicates

| Function | Algorithm | Evidence (`decimal.js` line) | Complexity | Dependencies | Go Plan |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `cmp` / `comparedTo` | Sign check $\rightarrow$ Exponent comparison $\rightarrow$ Digit-by-digit limb comparison | Lines 246-278 | $O(N)$ | Direct | Check signs; check exponents; loop through limb slices `xd` and `yd`; compare length if all limbs match. Return -1, 0, 1, or NaN. |
| `eq` / `equals` | Delegate to `cmp(y) == 0` | Lines 489-491 | $O(N)$ | `cmp` | Return `x.Cmp(y) == 0`. |
| `gt` | Delegate to `cmp(y) > 0` | Lines 509-511 | $O(N)$ | `cmp` | Return `x.Cmp(y) > 0`. |
| `gte` | Delegate to `cmp(y) >= 0` | Lines 519-522 | $O(N)$ | `cmp` | Return `x.Cmp(y) >= 0`. |
| `lt` | Delegate to `cmp(y) < 0` | Lines 925-927 | $O(N)$ | `cmp` | Return `x.Cmp(y) < 0`. |
| `lte` | Delegate to `cmp(y) <= 0` | Lines 934-937 | $O(N)$ | `cmp` | Return `x.Cmp(y) <= 0`. |
| `isNaN` | Check `s == 0` / `s == null` | Lines 889-891 | $O(1)$ | Direct | Return `d.Sign == 0` (or NaN sentinel). |
| `isFinite` | Check `d != nil && s != 0` | Lines 870-872 | $O(1)$ | Direct | Return `d.Coeff != nil && d.Sign != 0`. |
| `isInteger` | Check if fractional decimal places exist beyond exponent | Lines 879-882 | $O(N)$ | Direct | Inspect limb slice past `e` position; verify zero remainder. |
| `isNegative` | Check `s < 0` | Lines 898-900 | $O(1)$ | Direct | Return `d.Sign < 0`. |
| `isPositive` | Check `s > 0` | Lines 907-909 | $O(1)$ | Direct | Return `d.Sign > 0`. |
| `isZero` | Check `d[0] == 0` | Lines 916-918 | $O(1)$ | Direct | Return `d.Coeff[0] == 0`. |

---

## Rounding & Formatting

| Function | Algorithm | Evidence (`decimal.js` line) | Complexity | Dependencies | Go Plan |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `finalise` | Limb digit positioning $\rightarrow$ Rounding mode tie-breaking $\rightarrow$ Overflow/Underflow check $\rightarrow$ Trailing zero removal | Lines 2946-3110 | $O(N)$ | `digitsToString`, `getZeroString` | Determine rounding digit `rd` within target limb; evaluate 1 of 9 rounding modes; increment limb with carry if rounding up; check `minE`/`maxE`; trim trailing zero limbs. |
| `digitsToString` | Base $10^7$ limb string formatting with zero-padding | Lines 2520-2547 | $O(N)$ | `getZeroString` | Convert leading limb to string; format middle limbs with 7-digit zero padding (`%07d`); strip trailing zeros on final limb. |
| `finiteToString` | Format into fixed-point (`123.45`) or scientific (`1.2345e+2`) string | Lines 3113-3145 | $O(N)$ | `digitsToString` | Insert decimal point relative to exponent `e`; append `e+N` exponent suffix if `isExp` is true. |

---

## Power & Roots

| Function | Algorithm | Evidence (`decimal.js` line) | Complexity | Dependencies | Go Plan |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `sqrt` / `squareRoot` | Newton-Raphson iteration: $x_{k+1} = 0.5 \cdot (x_k + S / x_k)$ | Lines 1432-1522 | $O(N \cdot M \cdot \log P)$ | `divide`, `finalise`, `digitsToString` | Float64 initial seed; iterate Newton division until limb string convergence; check boundary nines/zeros. |
| `cbrt` / `cubeRoot` | Halley's method iteration: $r_{k+1} = r_k \cdot \frac{r_k^3 + 2X}{2r_k^3 + X}$ | Lines 334-421 | $O(N \cdot M \cdot \log P)$ | `divide`, `finalise`, `digitsToString` | Float64 initial seed; iterate Halley cubic division; check boundary rounding digits; finalize precision. |
| `pow` / `toPower` | Integer powers: Exponentiation by squaring (`intPow`). Non-integer: $x^y = \exp(y \cdot \ln(x))$ | Lines 2268-2365 | $O(M \cdot \text{Cost}(\exp, \ln))$ | `intPow`, `naturalExponential`, `naturalLogarithm`, `finalise` | If $y$ is small int, call binary exponentiation (`intPow`); else compute $\exp(y \times \ln(x))$ with guard digits. |

---

## Logarithms & Transcendentals

| Function | Algorithm | Evidence (`decimal.js` line) | Complexity | Dependencies | Go Plan |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `ln` / `naturalLogarithm` | Argument reduction ($x = m \cdot 10^k \rightarrow y = \frac{x-1}{x+1}$) $\rightarrow$ Series expansion $\ln(\frac{1+y}{1-y}) = 2 \sum \frac{y^{2k+1}}{2k+1}$ | Lines 3398-3512 | $O(P \cdot N^2)$ | `divide`, `finalise`, `getLn10` | Reduce argument near 1; compute series terms using quotient division; add $k \cdot \ln(10)$. |
| `exp` / `naturalExponential` | Argument reduction ($x = k \cdot \ln 2 + r$) $\rightarrow$ Taylor series $e^r = \sum \frac{r^n}{n!}$ | Lines 3307-3396 | $O(P \cdot N^2)$ | `divide`, `finalise`, `digitsToString` | Reduce argument by scaling by powers of 2; sum Taylor terms until term magnitude is below precision threshold. |
| `sin` / `cos` / `tan` | Argument reduction mod $2\pi$ to $[0, \pi/2]$ $\rightarrow$ Argument halving $\rightarrow$ Taylor series | Lines 2641-2672, 3687-3719, 3721-3755 | $O(P \cdot N^2)$ | `taylorSeries`, `toLessThanHalfPi`, `divide`, `finalise` | Reduce argument to $[0, \pi/4]$; evaluate Taylor series $\sum (-1)^n \frac{x^{2n}}{(2n)!}$; reconstruct via trig angle identities. |
