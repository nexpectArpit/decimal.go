# Remediation Plan

**Based on**: [FINAL_INDEPENDENT_VERIFICATION_AUDIT.md](FINAL_INDEPENDENT_VERIFICATION_AUDIT.md) (independent re-audit, score 41/100)
**Goal**: get from "looks done, isn't" to actually behaviorally equivalent with `decimal.js`, with an audit trail that can be trusted without re-verification.
**Principle**: fix the code first, harden the tests second, only then touch the score. No claim in any `.md` gets updated until the test that backs it has been re-run in this session.

---

## Phase 0 — Stop the bleeding (make the test suite tell the truth)

Nothing else in this plan is safe to act on until the tests are honest, because right now some "green" results are meaningless. Do this before any math fix.

### 0.1 Wire the real differential fuzzer into `go test`
- **Problem**: `fuzz/harness.go` (the one that caught `sin`/`atan`/`pow` returning garbage) is a `func main()`, never executed by `go test ./...`. Only `fuzz/differential_fuzz_test.go`'s `TestDifferentialFuzz` runs in CI, and it only checks `Add`/`Sub` against float64 at `1e-10` tolerance.
- **Fix**: convert `harness.go`'s logic into a `Test*` function inside `fuzz/differential_fuzz_test.go` (or a new `fuzz/transcendental_fuzz_test.go`). Cover at minimum: `Sin`, `Cos`, `Atan`, `Ln`, `Exp`, `Pow` (including negative-base/fractional-exponent). Fail the test (`t.Fatalf`) if divergence rate exceeds a defined threshold (start at 0%, relax only with a written justification per operation).
- **Acceptance**: `go test ./...` fails today, before any math fix, once this lands — that failure is the correct, honest state.
- **Effort**: S (half day).
- **Owner note**: do not delete `fuzz/log.txt` / `fuzz/harness.go`'s output — keep them as the historical record of what was found.

### 0.2 Stop the bridge from testing decimal.js against itself
- **Problem**: `tests/original/bridge.js` routes `dp`, `sd`/`precision`, `isInt`, no-arg `toFixed`, `toBinary`, `toHex`/`toHexadecimal`, `toOctal`, `toFraction`, `toNearest`, `clamp`/`clampedTo` through `toJS(this)` — i.e. the real npm `decimal.js`, not the Go CLI.
- **Fix options, in order of preference**:
  1. Implement the missing Go methods (`ToFraction`, `ToNearest`, `ToBinary`, `ToHex`, `ToOctal`, `Clamp`/`ClampedTo`, `Dp`, `Sd`, `IsInt`, no-arg `ToFixed`) in `src/`, expose them over the CLI's `--rpc` protocol, and repoint the bridge at `callGo(...)` like every other method.
  2. If a method is deliberately out of scope for this port, mark it `Skip` in the JS test runner (`tests/original/setup.js` / the relevant `modules/*.js`) with a comment explaining why, instead of silently faking a pass.
- **Acceptance**: `grep -n "toJS(this)" tests/original/bridge.js` returns zero results, OR every remaining hit has an adjacent `// INTENTIONALLY NOT PORTED:` comment and a matching `Skip` in the test runner.
- **Effort**: M (implementing ~6 new methods + RPC plumbing + bridge rewire — 2-3 days).

### 0.3 Delete or regenerate the fabricated benchmark file
- **Problem**: `bench/results.json` claims `fuzz_divergences: 0` and `memory_allocs_per_op: 1`, contradicted by `fuzz/log.txt` (19,129/19,538 divergences) and `PERFORMANCE_REPORT.md`'s own table (7-23 allocs/op, 5,259 for `Ln`).
- **Fix**: delete `bench/results.json`. Regenerate it (or replace it with a script, e.g. `make bench > bench/results.json`, that runs `go test -bench=. -benchmem ./tests/port/...` and captures real output) so it can never again be hand-edited into fiction.
- **Acceptance**: `bench/results.json` is either absent or produced by a committed, runnable script; its numbers match a benchmark run executed in the same session that updates it.
- **Effort**: XS (a few hours).

### 0.4 Re-verify `VERIFICATION_TABLE.json` row by row
- **Problem**: 2 of 2 checkable "VERIFIED_BUG" rows and 1 of 1 checkable "VERIFIED_MISSING" row did not reproduce against current `src/`.
- **Fix**: for every row, write a throwaway Go test reproducing the exact `input`/`precision` from the JSON, run it, and record the *actual* current output next to the claimed one. Delete or correct rows that don't reproduce; keep only what's freshly verified. Add a `verifiedAt` commit hash to each row going forward so this can't silently go stale again.
- **Acceptance**: every row's `actualGo` field matches a real, just-run test output; no row references code/line numbers that don't exist in the current tree.
- **Effort**: S (1 day, mostly mechanical).

### 0.5 Reconcile the failure-count discrepancy
- **Problem**: 410 (old audit) vs. 75 (`REMAINING_75_FAILURES_REPORT.md`) vs. 141 (actual current `FAILURE_DATABASE.csv`, properly CSV-parsed) — three numbers, no explanation.
- **Fix**: regenerate `FAILURE_DATABASE.csv` from a fresh run of `tests/original/collect_failures.js` against the current build, and update every doc that cites a failure count to reference that single fresh number with the date/commit it was generated at. Retire the old counts or clearly label them "historical, superseded."
- **Acceptance**: exactly one authoritative failure count exists at any time, timestamped and regenerable by one command.
- **Effort**: XS.

### 0.6 Explain the gap in `DECISIONS.md` (entries 16, 18, 20)
- **Fix**: check `git log`/`git blame` on `DECISIONS.md` for whether these entries were removed. If recoverable, restore them with a note on why they were superseded. If not recoverable, add a one-line note at each gap: "Decision N removed on [date]: [reason]" so the log is honest about not being contiguous by accident.
- **Effort**: XS.

---

## Phase 1 — Fix the catastrophic transcendental bugs

These are the correctness issues that actually matter to users of the library. Ordered by severity of observed divergence.

### 1.1 `Sin`/`Cos` blow up for `|x| > 2` (observed: `sin(-70.44)` → `9.99e+34` instead of `-0.97`)
- **Where to look**: `src/trig.go:86` (`Sin`), specifically `c.toLessThanHalfPi(x)` (range reduction) and the Taylor series loop at lines 114-131.
- **Hypothesis to test first**: `toLessThanHalfPi` is either not reducing large-magnitude inputs correctly (e.g. losing precision computing `x mod 2π` when `x`'s working precision `wpr` scales with `absE` but the `π` constant used for reduction doesn't have matching precision), or the quadrant tracking (`quadrant > 2` sign flip at line 133) is wrong for some quadrant, causing the series to run on an unreduced or wrongly-signed argument and diverge.
- **Verification step**: add a unit test that directly calls `toLessThanHalfPi(-70.438949)` and asserts the reduced value's magnitude is `< π/2` before touching the series code at all. If the reduced value is already wrong, the bug is upstream of the series; if it's right, the bug is in the series/sign logic.
- **Acceptance**: the Phase 0.1 fuzz test passes at 0% divergence (within tolerance) for `Sin`/`Cos` across the same random distribution that found this.

### 1.2 `Atan` blows up for `|x| > 1` (observed: `atan(2)` → `-3.2e+57`, `atan(5)` → `-6.0e+136`)
- **Where to look**: `src/trig.go:383` (`Atan`), the iterative halving reduction at lines 417-430 (`atan(x) = 2*atan(x/(1+sqrt(1+x^2)))` applied `k` times), and the final `2^k` rescale at lines 452-456.
- **Hypothesis to test first**: `k := wpr/LogBase + 2` (line 418) is a precision-derived iteration count, not one derived from how large `x` actually is — for `x=2` or `x=5` this reduction should shrink `xWork` toward 0 within a handful of iterations regardless of `wpr`, but if `wpr` is large, `k` is large, and repeatedly halving using `wprCtx.Sqrt`/`divide` at high working precision for ~28 iterations is exactly the kind of place compounding rounding or a wrong-precision `Sqrt` call could blow up rather than converge. Instrument `xWork` after each iteration of the loop at line 424 and confirm it's monotonically shrinking toward 0; if it isn't, the bug is in that loop, not the series.
- **Acceptance**: same fuzz test at 0% divergence for `Atan`, including the exact reported inputs `2.0` and `5.0` as permanent regression cases.

### 1.3 `Pow` with negative base + non-integer exponent must return `NaN`, not a ~1.5M-digit number
- **Observed**: `(-2)^0.5` → `9.5130527730906627269e+1505149` (should be `NaN`, since decimal.js and IEEE 754 both define this as undefined in the reals).
- **Where to look**: `src/pow.go`, the non-integer-exponent path of `toPower`/`Pow` (Decision 19 in `DECISIONS.md` only documents the negative-base *integer*-exponent case — the non-integer case has no corresponding decision or guard at all, which is the likely root cause: there's simply no negative-base check before falling into `exp(y * ln(x))`, and `ln` of a negative number is itself undefined, which is presumably where the runaway digit count originates).
- **Fix**: before computing `exp(y * ln(x))`, check `x.s < 0 && !y.isInteger()` and return `NaN` immediately, matching `decimal.js`'s `toPower` behavior.
- **Acceptance**: `Pow(-2, 0.5)` returns `NaN`; add this as a permanent unit test in `tests/port`, not just the fuzz suite (fuzz alone won't reliably hit exact half-integer exponents).

### 1.4 `Ln` accuracy drift for `x=1000` (observed: correct to 2 digits, wrong from the 3rd digit on — `6.599...` vs `6.9077...`, actually looks like more than a rounding issue, closer to a wrong argument-reduction constant)
- **Where to look**: `src/transcendental.go`, `Ln`'s argument-reduction step (Decision 7 references `y=(x-1)/(x+1)` reduction) and the static `ln(10)` constant used to rescale after reducing `x` into range.
- **Acceptance**: `Ln(1000)` matches `decimal.js` to full configured precision; add as a fixed regression test alongside the fuzz coverage.

### 1.5 Sweep for the same class of bug elsewhere
Given `Sin`, `Atan`, and `Pow` all failed catastrophically (not just last-digit) on inputs the curated `tests/original` suite didn't happen to hit, treat every transcendental (`Cos`, `Tan`, `Asin`, `Acos`, `Sinh`, `Cosh`, `Tanh`, `Asinh`, `Acosh`, `Atanh`, `Exp`, `Log`, `Log2`, `Sqrt`, `Cbrt`, `Hypot`) as suspect until the Phase 0.1 fuzz test has actually exercised each one across a wide input range (including large magnitude, near-boundary, and negative/special-value inputs) and shown 0% divergence.

---

## Phase 2 — Close the curated-suite gap (141 known failures)

Only start this after Phase 1, since some of these may turn out to be side effects of the same range-reduction bugs (e.g. `minus`/`plus` signed-zero failures are unlikely to be related, but `pow`/`sqrt`/`log`/`intPow` failures plausibly share root cause with Phase 1 items).

Priority order by current failure count (from `FAILURE_DATABASE.csv`, pending Phase 0.5 reconciliation):

1. **`plus` (84 failures)** — signed-zero (`+0` vs `-0`) and guard-digit boundary cases under `RoundFloor` for near-equal opposing-sign additions. `src/add.go`.
2. **`minus` (24 failures)** — same signed-zero issue, subtraction path. `src/sub.go`.
3. **`pow` (9), `sqrt` (9), `intPow` (7), `log` (4)** — re-run against the fuzz suite from Phase 1 first; fix any that are duplicates of the range-reduction bugs there before treating these as separate issues.
4. **`hypot`, `immutability`, `random`, `sign`** (1 each) — one-off edge cases, tackle last.

For each: reproduce the exact failing assertion from `FAILURE_DATABASE.csv` as a standalone Go unit test in `tests/port` *before* changing any code, so the fix is provably targeted.

---

## Phase 3 — Re-audit

Do not touch any score or "PASS" claim until:

1. `go test ./...` — including the new Phase 0.1 fuzz tests — passes with 0% divergence on transcendentals (or an explicitly justified nonzero tolerance).
2. `grep -n "toJS(this)" tests/original/bridge.js` is empty or fully justified (0.2).
3. `node tests/original/test.js` completes with the Phase 2 failure count at 0 (or a written, current explanation for any remainder).
4. `bench/results.json` is regenerated fresh (0.3).
5. `VERIFICATION_TABLE.json` and `DECISIONS.md` are internally consistent with the current tree (0.4, 0.6).

Then, and only then, write a new independent audit — do not edit the 41/100 score in place; append a dated re-audit the same way this one superseded the 84/100 one, so the score history stays honest.

---

## Effort Summary

| Phase | Item | Effort |
|---|---|---|
| 0.1 | Wire fuzzer into `go test` | S |
| 0.2 | Fix bridge fallthroughs | M |
| 0.3 | Fix benchmark file | XS |
| 0.4 | Re-verify VERIFICATION_TABLE.json | S |
| 0.5 | Reconcile failure counts | XS |
| 0.6 | Explain DECISIONS.md gap | XS |
| 1.1 | Fix Sin/Cos range reduction | M-L |
| 1.2 | Fix Atan reduction loop | M |
| 1.3 | Fix Pow negative-base/fractional-exponent | S |
| 1.4 | Fix Ln(1000)-class drift | M |
| 1.5 | Fuzz sweep remaining transcendentals | M (mostly Phase 0.1 reuse) |
| 2 | Close 141 curated failures | M-L |
| 3 | Re-audit | S |

(XS = hours, S = ~1 day, M = 2-3 days, L = a week+)

**Do Phase 0 first, in full, before writing a single line of math fix.** Otherwise there is no way to know whether a "fix" actually fixed anything, or just moved the goalposts the way the current docs did.
