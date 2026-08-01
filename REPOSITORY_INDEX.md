# Repository Index: decimal.js v10.6.0

## Status: COMPLETE

This document provides a comprehensive index of all constants, internal helper functions, prototype methods, and static methods in `decimal.js`, mapped to exact line numbers and dependency graphs.

---

## Constants & Configuration

| Name | Value | Location in `decimal.js` | Purpose | Risk / Go Translation Notes |
| :--- | :--- | :--- | :--- | :--- |
| `EXP_LIMIT` | `9e15` | Line 19 | Exponent magnitude limit for `minE`, `maxE`, `toExpNeg`, `toExpPos` | Represented via Go `int64` |
| `MAX_DIGITS` | `1e9` | Line 23 | Maximum precision limit | Checked in config validator |
| `NUMERALS` | `'0123456789abcdef'` | Line 26 | Base conversion alphabet | String lookup or slice index |
| `LN10` | `2.3025...` (1025 digits) | Line 28 | High-precision constant $\ln(10)$ | Pre-computed static string |
| `PI` | `3.1415...` (1025 digits) | Line 31 | High-precision constant $\pi$ | Pre-computed static string |
| `DEFAULTS` | `{ precision: 20, rounding: 4, ... }` | Lines 36-95 | Default global configuration settings | Package defaults & context initializer |
| `BASE` | `1e7` ($10^7$) | Line 118 | Radix base for limb coefficient array `d` | `const Base int32 = 10000000` |
| `LOG_BASE` | `7` | Line 119 | Number of decimal digits per limb | `const LogBase int = 7` |
| `MAX_SAFE_INTEGER` | `9007199254740991` | Line 120 | $2^{53}-1$ boundary for exact JS numbers | Native Go `int64` exceeds this |

---

## Internal Helper Functions

| Name | Line Numbers | Purpose | Called By | Calls | Tests Validating | Priority |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `finalise` | 2946-3110 | Rounds coefficient array, checks minE/maxE, trims zeros | Almost ALL math methods | `digitsToString`, `getZeroString` | `test/modules/*.js` | **CRITICAL (1)** |
| `divide` | 2728-2939 | Base $10^7$ long division with quotient estimation & normalization | `div`, `divToInt`, `mod`, `sqrt`, `cbrt`, `toFraction`, `tan`, `tanh`, transcendentals | `multiplyInteger`, `compare`, `subtract`, `finalise` | `test/modules/dividedBy.js` | **CRITICAL (2)** |
| `digitsToString` | 2520-2547 | Converts coefficient array `d` to decimal string | `cbrt`, `toFraction`, `pow`, `finiteToString`, `naturalExponential`, `naturalLogarithm` | `getZeroString` | `test/modules/toString.js` | **HIGH (3)** |
| `parseDecimal` | 3523-3605 | Parses standard decimal and scientific notation strings | Constructor (`Decimal`) | `parseDecimal` | `test/modules/constructor.js` | **HIGH (4)** |
| `parseOther` | 3607-3685 | Parses binary (`0b`), octal (`0o`), hex (`0x`) strings | Constructor (`Decimal`) | `convertBase`, `divide` | `test/modules/constructor.js` | **HIGH (5)** |
| `checkInt32` | 2550-2554 | Validates integer parameter boundaries | `toDP`, `toFixed`, `toPrecision`, `toSD`, `toNearest` | None | `test/modules/*.js` | **MEDIUM (6)** |
| `checkRoundingDigits` | 2562-2607 | Checks 4/5 rounding digits near boundaries | `pow`, `naturalExponential`, `naturalLogarithm` | None | `test/modules/toPower.js` | **MEDIUM (7)** |
| `convertBase` | 2613-2634 | Converts digit string between arbitrary bases | `parseOther`, `toStringBinary` | None | `test/modules/toHex.js` | **MEDIUM (8)** |
| `finiteToString` | 3113-3145 | Converts finite Decimal to string format | `toString`, `toFixed`, `toExponential`, `toPrecision`, `valueOf` | `digitsToString`, `getZeroString` | `test/modules/toString.js` | **HIGH (9)** |
| `getBase10Exponent` | 3147-3154 | Calculates base-10 exponent from digit array | `plus`, `minus`, `times`, `parseOther` | None | `test/modules/plus.js` | **MEDIUM (10)** |
| `getLn10` | 3156-3166 | Returns or computes $\ln(10)$ to requested precision | `logarithm`, `naturalLogarithm` | `finalise` | `test/modules/ln.js` | **MEDIUM (11)** |
| `getPi` | 3168-3172 | Returns or computes $\pi$ to requested precision | `acos`, `asin`, `atan`, `atan2` | `finalise` | `test/modules/cos.js` | **MEDIUM (12)** |
| `getPrecision` | 3174-3192 | Calculates significant digits of digit array `d` | `precision`/`sd`, `toFraction` | None | `test/modules/precision.js` | **MEDIUM (13)** |
| `intPow` | 3208-3241 | Exponentiation by squaring for integer powers | `pow`, `parseOther` | `truncate` | `test/modules/toPower.js` | **HIGH (14)** |
| `naturalExponential`| 3307-3396 | Computes $e^x$ via argument reduction & series | `exp`, `pow` | `divide`, `finalise`, `digitsToString` | `test/modules/exp.js` | **HIGH (15)** |
| `naturalLogarithm`  | 3398-3512 | Computes $\ln(x)$ via series expansion | `ln`, `log`, `acosh`, `asinh`, `atanh`, `pow` | `divide`, `finalise`, `getLn10` | `test/modules/ln.js` | **HIGH (16)** |
| `taylorSeries` | 3721-3755 | Taylor series evaluator for trig & hyperbolic functions | `cos`, `cosh`, `sin`, `sinh`, `cosine`, `sine` | `divide` | `test/modules/cos.js` | **HIGH (17)** |
| `toStringBinary` | 3803-3934 | Formats Decimal into hex, binary, or octal string | `toHex`, `toBinary`, `toOctal` | `convertBase`, `divide`, `finiteToString` | `test/modules/toHex.js` | **MEDIUM (18)** |

---

## Prototype Methods (Public API Surface)

| Method Name | Alias | Location in `decimal.js` | Core Helper Invoked | Purpose | Priority |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `absoluteValue` | `abs` | 200-204 | `finalise` | Absolute value | High |
| `ceil` | - | 212-214 | `finalise` | Round towards +Infinity | High |
| `clampedTo` | `clamp` | 225-235 | `cmp`, Constructor | Clamp value in range [min, max] | Medium |
| `comparedTo` | `cmp` | 246-278 | Direct comparison | Compare two Decimals (-1, 0, 1, NaN) | Critical |
| `cosine` | `cos` | 294-315 | `cosine`, `toLessThanHalfPi`, `finalise` | Cosine in radians | High |
| `cubeRoot` | `cbrt` | 334-421 | `divide`, `finalise`, `digitsToString` | Cube root | High |
| `decimalPlaces` | `dp` | 428-444 | Direct calculation | Count fractional decimal places | High |
| `dividedBy` | `div` | 468-470 | `divide` | Division | Critical |
| `dividedToIntegerBy` | `divToInt` | 478-482 | `divide`, `finalise` | Integer division | High |
| `equals` | `eq` | 489-491 | `cmp` | Value equality check | Critical |
| `floor` | - | 499-501 | `finalise` | Round towards -Infinity | High |
| `greaterThan` | `gt` | 509-511 | `cmp` | Greater than check | Critical |
| `greaterThanOrEqualTo`| `gte` | 519-522 | `cmp` | Greater than or equal check | Critical |
| `hyperbolicCosine` | `cosh` | 550-590 | `taylorSeries`, `finalise` | Hyperbolic cosine | Medium |
| `hyperbolicSine` | `sinh` | 601-638 | `taylorSeries`, `finalise` | Hyperbolic sine | Medium |
| `hyperbolicTangent` | `tanh` | 647-669 | `divide`, `finalise` | Hyperbolic tangent | Medium |
| `inverseCosine` | `acos` | 679-699 | `getPi`, `asin`, `finalise` | Inverse cosine (acos) | Medium |
| `inverseHyperbolicCosine` | `acosh` | 709-725 | `naturalLogarithm`, `sqrt`, `finalise` | Inverse hyperbolic cosine | Medium |
| `inverseHyperbolicSine` | `asinh` | 735-752 | `naturalLogarithm`, `sqrt`, `finalise` | Inverse hyperbolic sine | Medium |
| `inverseHyperbolicTangent` | `atanh` | 762-780 | `naturalLogarithm`, `divide`, `finalise` | Inverse hyperbolic tangent | Medium |
| `inverseSine` | `asin` | 790-820 | `getPi`, `atan`, `sqrt`, `finalise` | Inverse sine (asin) | Medium |
| `inverseTangent` | `atan` | 830-863 | `getPi`, `finalise` | Inverse tangent (atan) | Medium |
| `isFinite` | - | 870-872 | Direct check | Check if finite | High |
| `isInteger` | `isInt` | 879-882 | Direct check | Check if integer | High |
| `isNaN` | - | 889-891 | Direct check | Check if NaN | High |
| `isNegative` | `isNeg` | 898-900 | Direct check | Check if negative | High |
| `isPositive` | `isPos` | 907-909 | Direct check | Check if positive | High |
| `isZero` | - | 916-918 | Direct check | Check if zero | High |
| `lessThan` | `lt` | 925-927 | `cmp` | Less than check | Critical |
| `lessThanOrEqualTo` | `lte` | 934-937 | `cmp` | Less than or equal check | Critical |
| `logarithm` | `log` | 954-1002 | `naturalLogarithm`, `getLn10`, `finalise` | Logarithm in base N | High |
| `minus` | `sub` | 1024-1110 | Direct subtraction, `finalise` | Subtraction | Critical |
| `modulo` | `mod` | 1133-1185 | `divide`, `finalise` | Modulo | Critical |
| `naturalExponential` | `exp` | 1195-1209 | `naturalExponential`, `finalise` | Natural exponential ($e^x$) | High |
| `naturalLogarithm` | `ln` | 1219-1229 | `naturalLogarithm`, `finalise` | Natural logarithm ($\ln(x)$) | High |
| `negated` | `neg` | 1236-1240 | `finalise` | Negation (-x) | High |
| `plus` | `add` | 1262-1358 | Direct addition, `finalise` | Addition | Critical |
| `precision` | `sd` | 1368-1377 | `getPrecision` | Get significant digits count | High |
| `round` | - | 1386-1388 | `finalise` | Round to whole number | High |
| `sine` | `sin` | 1398-1419 | `sine`, `toLessThanHalfPi`, `finalise` | Sine in radians | High |
| `squareRoot` | `sqrt` | 1432-1522 | `divide`, `finalise`, `digitsToString` | Square root | High |
| `tangent` | `tan` | 1533-1555 | `sine`, `cosine`, `divide`, `finalise` | Tangent in radians | High |
| `times` | `mul` | 1577-1678 | Direct multiplication, `finalise` | Multiplication | Critical |
| `toBinary` | - | 1690-1692 | `toStringBinary` | String in binary | Medium |
| `toDecimalPlaces` | `toDP` | 1704-1718 | `finalise` | Round to N decimal places | High |
| `toExponential` | - | 1730-1746 | `finalise`, `finiteToString` | String in exponential notation | High |
| `toFixed` | - | 2026-2046 | `finalise`, `finiteToString` | String in fixed-point notation | High |
| `toFraction` | - | 2060-2118 | `divide`, `finalise`, `digitsToString` | Rational fraction representation | Medium |
| `toHexadecimal` | `toHex` | 2131-2133 | `toStringBinary` | String in hexadecimal | Medium |
| `toNearest` | - | 2152-2197 | `divide`, `finalise` | Round to nearest multiple | Medium |
| `toNumber` | - | 2205-2207 | Direct double conversion | Convert to float64 number | High |
| `toOctal` | - | 2220-2222 | `toStringBinary` | String in octal | Medium |
| `toPower` | `pow` | 2268-2365 | `intPow`, `naturalExponential`, `naturalLogarithm`, `finalise` | Power $x^y$ | High |
| `toPrecision` | - | 2379-2397 | `finalise`, `finiteToString` | String to N precision | High |
| `toSignificantDigits`| `toSD` | 2414-2429 | `finalise` | Round to N significant digits | High |
| `toString` | - | 2439-2445 | `finiteToString` | Convert Decimal to string | Critical |
| `truncated` | `trunc` | 2452-2454 | `finalise` | Truncate fraction | High |
| `valueOf` | `toJSON` | 2462-2468 | `finiteToString` | JSON / value primitive string | High |
