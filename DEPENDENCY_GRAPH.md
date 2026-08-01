# Dependency Graph & Bottom-Up Implementation Order

## Status: COMPLETE (Derived from `decimal.js` Internal Call Tree)

This document establishes the exact internal call dependencies of `decimal.js` and derives the strict bottom-up implementation ordering for `our-projectInGO`.

---

## Internal Call Tree

```
[Level 0: Foundation Primitives]
├── LOG_BASE (7), BASE (1e7), EXP_LIMIT (9e15)
├── getZeroString()
├── checkInt32()
└── isOdd()

[Level 1: String Formatting & Base Utility]
├── digitsToString() ────────► getZeroString()
├── convertBase()
├── getBase10Exponent()
└── checkRoundingDigits()

[Level 2: Core Alignment & Normalization]
├── finalise() ──────────────► digitsToString(), getZeroString()
├── compare()
├── multiplyInteger()
└── subtract()

[Level 3: Core Arithmetic & Division Engine]
├── divide() ────────────────► multiplyInteger(), compare(), subtract(), finalise()
├── add() ───────────────────► getBase10Exponent(), finalise()
├── sub() ───────────────────► getBase10Exponent(), finalise()
└── mul() ───────────────────► finalise()

[Level 4: Decimal Constructors & Parsers]
├── parseDecimal() ──────────► Decimal constructor
├── parseOther() ────────────► convertBase(), divide()
└── Decimal() ───────────────► parseDecimal(), parseOther()

[Level 5: Prototype Relational & Basic Arithmetic Wrappers]
├── cmp() ───────────────────► Direct comparison
├── eq(), gt(), gte(), lt(), lte() ──► cmp()
├── plus() ──────────────────► add()
├── minus() ─────────────────► sub()
├── times() ─────────────────► mul()
├── dividedBy() ─────────────► divide()
├── dividedToIntegerBy() ────► divide(), finalise()
└── modulo() ────────────────► divide(), finalise()

[Level 6: Transcendental Base Helpers]
├── getLn10() ───────────────► finalise()
├── getPi() ─────────────────► finalise()
├── naturalLogarithm() ──────► divide(), finalise(), getLn10()
├── naturalExponential() ────► divide(), finalise(), digitsToString()
├── taylorSeries() ──────────► divide()
└── intPow() ────────────────► truncate()

[Level 7: Advanced Mathematical Operations]
├── sqrt() ──────────────────► divide(), finalise(), digitsToString()
├── cbrt() ──────────────────► divide(), finalise(), digitsToString()
├── pow() ───────────────────► intPow(), naturalExponential(), naturalLogarithm(), finalise()
├── sin(), cos(), tan() ─────► taylorSeries(), toLessThanHalfPi(), divide(), finalise()
├── asin(), acos(), atan() ──► getPi(), naturalLogarithm(), sqrt(), finalise()
└── sinh(), cosh(), tanh() ──► taylorSeries(), divide(), finalise()

[Level 8: Formatting & Public Conversions]
├── finiteToString() ────────► digitsToString(), getZeroString()
├── toStringBinary() ────────► convertBase(), divide(), finiteToString()
├── toString() ──────────────► finiteToString()
├── toFixed() ───────────────► finalise(), finiteToString()
├── toExponential() ─────────► finalise(), finiteToString()
├── toPrecision() ───────────► finalise(), finiteToString()
├── toFraction() ────────────► divide(), finalise(), digitsToString()
└── toHex(), toBinary(), toOctal() ──► toStringBinary()
```

---

## Strict Bottom-Up Implementation Order for Go

To ensure zero forward-declaration stubs or unfulfilled dependencies during development, Go files will be authored in this exact sequence:

1. **`src/types.go`** (Constants, `BASE=1e7`, `Decimal` struct, `Context` struct, Error types)
2. **`src/utils.go`** (`getZeroString`, `checkInt32`, `isOdd`, `digitsToString`, `convertBase`)
3. **`src/finalise.go`** (`finalise` rounding engine, 9 rounding modes, overflow/underflow checks)
4. **`src/compare.go`** (`compare`, `cmp`, `eq`, `gt`, `gte`, `lt`, `lte`, predicate methods `isNaN`, `isFinite`, `isZero`, `isNeg`, `isPos`, `isInt`)
5. **`src/add.go` & `src/sub.go`** (`getBase10Exponent`, `add`/`plus`, `sub`/`minus`, `abs`, `neg`)
6. **`src/mul.go`** (`multiplyInteger`, `mul`/`times`)
7. **`src/div.go` & `src/mod.go`** (`divide`, `div`/`dividedBy`, `divToInt`, `mod`/`modulo`)
8. **`src/parser.go` & `src/constructor.go`** (`parseDecimal`, `parseOther`, `New`, `clone`, `config`)
9. **`src/pow.go` & `src/roots.go`** (`intPow`, `sqrt`, `cbrt`, `pow`)
10. **`src/transcendental.go`** (`getLn10`, `getPi`, `naturalLogarithm`/`ln`, `naturalExponential`/`exp`, `log`)
11. **`src/trig.go` & `src/hyperbolic.go`** (`taylorSeries`, `toLessThanHalfPi`, `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`, `sinh`, `cosh`, `tanh`, `asinh`, `acosh`, `atanh`)
12. **`src/format.go` & `src/convert.go`** (`finiteToString`, `toStringBinary`, `toString`, `toFixed`, `toExponential`, `toPrecision`, `toNumber`, `toFraction`, `toHex`, `toBinary`, `toOctal`, `toJSON`)
