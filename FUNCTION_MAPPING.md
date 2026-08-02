# FUNCTION MAPPING — decimal.js → decimal.go

**Status**: IN_PROGRESS

## Core Arithmetic
| decimal.js Method | JS Line | Go Function | Go File:Line | Status |
|---|---|---|---|---|
| `P.plus / P.add` | 1262 | `Context.Add`, `Decimal.Add`, `Decimal.Plus` | add.go:9 | ⚠️ NEEDS REVIEW (H1: no external flag) |
| `P.minus / P.sub` | 1264 | `Context.Sub`, `Decimal.Sub`, `Decimal.Minus` | sub.go:9 | ⚠️ NEEDS REVIEW |
| `P.times / P.mul` | 1577 | `Context.Mul`, `Decimal.Mul`, `Decimal.Times` | mul.go:35 | ⚠️ NEEDS REVIEW |
| `P.dividedBy / P.div` | 468 | `Context.Div`, `Decimal.Div`, `Decimal.DividedBy` | div.go:172 | ⚠️ NEEDS REVIEW |
| `P.modulo / P.mod` | 1133 | `Context.Mod`, `Decimal.Mod`, `Decimal.Modulo` | mod.go:5 | ❌ BROKEN (H4: ignores modulo mode) |
| `P.dividedToIntegerBy / P.divToInt` | 478 | `Context.DivToInt`, `Decimal.DivToInt` | mod.go:36 | ⚠️ NEEDS REVIEW |

## Comparison
| decimal.js Method | JS Line | Go Function | Go File:Line | Status |
|---|---|---|---|---|
| `P.comparedTo / P.cmp` | 246 | `Decimal.Cmp` | compare.go:27 | ❌ BROKEN (H7) |
| `P.equals / P.eq` | 486 | `Decimal.Eq` | compare.go:111 | ⚠️ Depends on Cmp |
| `P.greaterThan / P.gt` | 504 | `Decimal.Gt` | compare.go:116 | ⚠️ Depends on Cmp |
| `P.greaterThanOrEqualTo / P.gte` | 514 | `Decimal.Gte` | compare.go:121 | ⚠️ Depends on Cmp |
| `P.lessThan / P.lt` | 524 | `Decimal.Lt` | compare.go:127 | ⚠️ Depends on Cmp |
| `P.lessThanOrEqualTo / P.lte` | 534 | `Decimal.Lte` | compare.go:132 | ⚠️ Depends on Cmp |
| `P.isNaN` | 694 | `Decimal.IsNaN` | compare.go:138 | ✅ OK |
| `P.isFinite` | 686 | `Decimal.IsFinite` | compare.go:143 | ✅ OK |
| `P.isZero` | 710 | `Decimal.IsZero` | compare.go:148 | ✅ OK |
| `P.isNeg / P.isNegative` | 694 | `Decimal.IsNeg` | compare.go:153 | ✅ OK |
| `P.isPos / P.isPositive` | 702 | `Decimal.IsPos` | compare.go:158 | ✅ OK |
| `P.isInt / P.isInteger` | 678 | `Decimal.IsInt` | compare.go:163 | ⚠️ M1 |

## Formatting
| decimal.js Method | JS Line | Go Function | Go File:Line | Status |
|---|---|---|---|---|
| `P.toString` | 2439 | `Decimal.String` | format.go:75 | ⚠️ NEEDS REVIEW |
| `P.valueOf` | 2449 | `Decimal.ValueOf` | format.go:184 | ⚠️ H5 |
| `P.toFixed` | 2026 | `Context.ToFixed`, `Decimal.ToFixed` | format.go:102 | ⚠️ NEEDS REVIEW |
| `P.toExponential` | 1730 | `Context.ToExponential`, `Decimal.ToExponential` | format.go:123 | ⚠️ NEEDS REVIEW |
| `P.toPrecision` | 2379 | `Context.ToPrecision`, `Decimal.ToPrecision` | format.go:144 | ⚠️ NEEDS REVIEW |
| `P.toNumber` | 2205 | `Decimal.Float64`, `Decimal.ToNumber` | convert.go:10 | ❌ BROKEN (H8) |
| `P.toJSON` | 2179 | `Decimal.MarshalJSON` | format.go:203 | ⚠️ NEEDS REVIEW |

## Transcendental / Advanced
| decimal.js Method | JS Line | Go Function | Go File:Line | Status |
|---|---|---|---|---|
| `P.naturalLogarithm / P.ln` | 1108 | `Context.Ln`, `Decimal.Ln` | transcendental.go:25 | ❌ CRITICAL (C3) |
| `P.naturalExponential / P.exp` | 502 | `Context.Exp`, `Decimal.Exp` | transcendental.go:78 | ❌ CRITICAL (C4) |
| `P.logarithm / P.log` | 1090 | `Context.Log`, `Decimal.Log` | transcendental.go:117 | ❌ BROKEN (depends C3) |
| `P.toPower / P.pow` | 2268 | `Context.Pow`, `Decimal.Pow`, `Decimal.ToPower` | pow.go:27 | ❌ CRITICAL (C5, C6) |
| `P.squareRoot / P.sqrt` | 1726 | `Context.Sqrt`, `Decimal.Sqrt`, `Decimal.SquareRoot` | roots.go:5 | ❌ CRITICAL (C9) |
| `P.cubeRoot / P.cbrt` | 334 | `Context.Cbrt`, `Decimal.Cbrt`, `Decimal.CubeRoot` | roots.go:48 | ⚠️ H6 |

## Trigonometric
| decimal.js Method | JS Line | Go Function | Go File:Line | Status |
|---|---|---|---|---|
| `P.sine / P.sin` | 1692 | `Context.Sin`, `Decimal.Sin`, `Decimal.Sine` | trig.go:5 | ❌ CRITICAL (C7) |
| `P.cosine / P.cos` | 294 | `Context.Cos`, `Decimal.Cos`, `Decimal.Cosine` | trig.go:54 | ❌ CRITICAL (C8) |
| `P.tangent / P.tan` | 1815 | `Context.Tan`, `Decimal.Tan`, `Decimal.Tangent` | trig.go:100 | ❌ BROKEN (depends C7, C8) |
| `P.inverseSine / P.asin` | 122 | `Context.Asin`, `Decimal.Asin` | trig.go:117 | ❌ CRITICAL (C11) |
| `P.inverseCosine / P.acos` | 87 | `Context.Acos`, `Decimal.Acos` | trig.go:133 | ❌ BROKEN (depends C11) |
| `P.inverseTangent / P.atan` | 159 | `Context.Atan`, `Decimal.Atan` | trig.go:148 | ❌ CRITICAL (C10) |
| `P.hyperbolicSine / P.sinh` | 583 | `Context.Sinh`, `Decimal.Sinh` | hyperbolic.go:4 | ❌ BROKEN (depends C4) |
| `P.hyperbolicCosine / P.cosh` | 422 | `Context.Cosh`, `Decimal.Cosh` | hyperbolic.go:31 | ❌ BROKEN (depends C4) |
| `P.hyperbolicTangent / P.tanh` | 1851 | `Context.Tanh`, `Decimal.Tanh` | hyperbolic.go:58 | ❌ BROKEN (depends C4) |
| `P.inverseHyperbolicSine / P.asinh` | 149 | `Context.Asinh`, `Decimal.Asinh` | hyperbolic.go:75 | ❌ BROKEN (depends C3) |
| `P.inverseHyperbolicCosine / P.acosh` | 76 | `Context.Acosh`, `Decimal.Acosh` | hyperbolic.go:90 | ❌ BROKEN (depends C3) |
| `P.inverseHyperbolicTangent / P.atanh` | 176 | `Context.Atanh`, `Decimal.Atanh` | hyperbolic.go:105 | ❌ BROKEN (depends C3) |

## Rounding
| decimal.js Method | JS Line | Go Function | Go File:Line | Status |
|---|---|---|---|---|
| `P.truncated / P.trunc` | 2430 | `Context.Trunc`, `Decimal.Trunc` | finalise.go:192 | ⚠️ NEEDS REVIEW |
| `P.floor` | 498 | `Context.Floor`, `Decimal.Floor` | finalise.go:205 | ⚠️ NEEDS REVIEW |
| `P.ceil` | 236 | `Context.Ceil`, `Decimal.Ceil` | finalise.go:218 | ⚠️ NEEDS REVIEW |
| `P.round` | 1358 | `Context.Round`, `Decimal.Round` | finalise.go:231 | ⚠️ NEEDS REVIEW |
| `P.absoluteValue / P.abs` | 73 | `Decimal.Abs` | sub.go:162 | ⚠️ NEEDS REVIEW |
| `P.negated / P.neg` | 1190 | `Decimal.Neg` | sub.go:174 | ⚠️ NEEDS REVIEW |

## MISSING from Go (NOT IMPLEMENTED)
| decimal.js Method | JS Line | Status |
|---|---|---|
| `P.toDecimalPlaces / P.toDP` | 1986 | ❌ MISSING |
| `P.toSignificantDigits / P.toSD` | 2397 | ❌ MISSING |
| `P.precision / P.sd` | 1343 | ❌ MISSING |
| `P.decimalPlaces / P.dp` | 460 | ❌ MISSING |
| `P.toFraction` | 2062 | ❌ MISSING |
| `P.toNearest` | 2121 | ❌ MISSING |
| `P.toBinary` | 1898 | ❌ MISSING |
| `P.toHexadecimal / P.toHex` | 2152 | ❌ MISSING |
| `P.toOctal` | 2198 | ❌ MISSING |
| `Decimal.atan2` | 4399 | ❌ MISSING |
| `Decimal.max` | 4510 | ❌ MISSING |
| `Decimal.min` | 4520 | ❌ MISSING |
| `Decimal.random` | 4565 | ❌ MISSING |
| `Decimal.sign` | 4617 | ❌ MISSING |
| `Decimal.sum` | 4630 | ❌ MISSING |
| `Decimal.clamp` | 4648 | ❌ MISSING |
| `P.clampedTo / P.clamp` | ~240 | ❌ MISSING |

## Summary
- **Total decimal.js methods**: ~57
- **Implemented in Go**: ~40
- **Working correctly**: ~12 (predicates + basic formatting)
- **Broken or wrong algorithm**: ~18
- **Missing entirely**: ~17
