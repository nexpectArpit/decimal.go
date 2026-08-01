# Project Agenda & Engineering Charter

## Mission Objective

Build a production-quality Go implementation of `decimal.js` that exhibits **100% behavioral equivalence** to the original JavaScript library, with **zero test modification** and **zero hacky workarounds**.

**Target Directory:** `/Users/arpittripathi/Desktop/coderescurration-project/our-projectInGO`

---

## Core Mindset

> **"Every time you think 'I know how this should be implemented,' stop and verify whether MikeMcl already solved the problem differently."**

We are NOT designing a decimal library.

We are reconstructing MikeMcl's engineering decisions in Go.

---

## 1. Track Assignment

- **Track:** `F` — **Runtime Modernization (JavaScript → Go)**
- **Difficulty:** Easy-Med | **Team:** 1-2 | **Pool:** ~6 repos | **Target LOC:** 2k - 5k
- **Mission Goal:** Eliminate the Node runtime dependency entirely.

---

## 2. Project Phases (Strict Sequential Order)

### Phase A: Research (NO CODE)

| Step | Document | Purpose | Gate |
| :--- | :--- | :--- | :--- |
| A1 | `REPOSITORY_INDEX.md` | Index every function, helper, constant in decimal.js | Must be 100% complete |
| A2 | `RESEARCH_PROTOCOL.md` | Answer every reverse-engineering question with evidence | Every `?` answered |
| A3 | `COMPARATIVE_ANALYSIS.md` | Compare solutions across all reference repos | Every `?` answered |
| A4 | `AUTHOR_PHILOSOPHY.md` | Understand why MikeMcl made each choice | Evidence-backed |

**Until Phase A is complete, no Go code is written. No exceptions.**

### Phase B: Design Decisions

| Step | Document | Purpose | Gate |
| :--- | :--- | :--- | :--- |
| B1 | `DECISIONS.md` | Record each architectural decision with evidence | Every decision justified |

Design decisions are derived FROM research, not invented before it.

### Phase C: Implementation

| Step | Action | Gate |
| :--- | :--- | :--- |
| C1 | Implement module-by-module, leaf functions first | Each module compiles |
| C2 | Run original decimal.js tests against Go implementation | Tests pass |
| C3 | Fix failures by investigating behavior, not adjusting tests | Behavior matches |

### Phase D: Validation

| Step | Action | Gate |
| :--- | :--- | :--- |
| D1 | Differential fuzzing (60s+ with 0 divergences) | `fuzz/log.txt` clean |
| D2 | Benchmarking (p99, RSS, startup, throughput) | `bench/results.json` honest |
| D3 | Verify 0 `unsafe` count | Static analysis pass |

### Phase E: Demo (User Manual Work)

5-minute live demo video — handled by team lead post-engineering.

---

## 3. Hackathon Scoring Criteria

| Criterion | Weight | What Judges Evaluate |
| :--- | :---: | :--- |
| **Functionality & Reliability** | **40%** | One-command build. Original tests pass unmodified. File-hash verified. |
| **Behavioral Equivalence** | **30%** | Differential fuzz survival. Honest p99/RSS/startup benchmarks with methodology. |
| **Code Quality** | **20%** | Idiomatic Go. Zero `unsafe`/`any` escape hatches. Quality of `DECISIONS.md`. |
| **Innovation** | **10%** | Latent bugs caught via differential testing. Architectural choices worth adopting upstream. |

---

## 4. Bonus Points (+16 Max)

| Category | Points | Difficulty | Requirement |
| :--- | :---: | :---: | :--- |
| **Differential Fuzz Survivor** | +5 | HARD | 60s+ fuzz, 0 divergences, published log |
| **Zero Unsafe** | +5 | HARD | 0 `unsafe`/`any`/escape-hatch count |
| **Decision Log** | +3 | MEDIUM | 10+ non-trivial decisions with rationale |
| **Bug Catcher** | +3 | MEDIUM | Discover and document latent bug in original |

---

## 5. Prohibited Anti-Patterns

| # | Anti-Pattern | Our Guarantee |
| :---: | :--- | :--- |
| ❌ 1 | Hello-world / single-function rewrites | Full module-by-module port |
| ❌ 2 | Shelling out to Node binary | 100% native Go |
| ❌ 3 | FFI / CGO into JS runtime | Pure Go, no CGO |
| ❌ 4 | Silently editing original tests | Tests remain untouched |
| ❌ 5 | Cherry-picking happy-path tests only | Full edge-case coverage |
| ❌ 6 | Unexplained AI code dumps | Every decision documented and defensible |
| ❌ 7 | Repos over 8,000 source lines | Strict 2k-5k LOC scope |
| ❌ 8 | Custom hardware / GUI / proprietary toolchains | Standard Go toolchain only |

---

## 6. Source of Truth & Reference Repositories

| Repository | Path | Role |
| :--- | :--- | :--- |
| **decimal.js** | `decimal.js/` | **SPECIFICATION.** Defines all behavior. |
| **bignumber.js** | `bignumber.js/` | **Author philosophy.** Understand MikeMcl's thinking. |
| **shopspring/decimal** | `shopspring-decimal/` | **Go engineering reference.** Package layout, idioms. |
| **cockroachdb/apd** | `apd/` | **Algorithm reference.** Study approaches to sqrt, ln, exp, etc. |
| **Go stdlib** | `math/big` | **Primitives reference.** Large integer mechanics. |

**Only `decimal.js` defines expected behavior. Everything else is educational.**

---

## 7. Submission Structure

```
our-projectInGO/
├── README.md
├── DECISIONS.md
├── RESEARCH_PROTOCOL.md
├── REPOSITORY_INDEX.md
├── COMPARATIVE_ANALYSIS.md
├── AUTHOR_PHILOSOPHY.md
├── Dockerfile
├── Makefile
├── .port-mortem.toml
├── src/
├── tests/
│   ├── original/
│   └── port/
├── fuzz/
│   ├── harness.go
│   └── log.txt
└── bench/
    ├── methodology.md
    └── results.json
```
