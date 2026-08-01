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
