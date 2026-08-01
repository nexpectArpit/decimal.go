# Engineering Decisions Log (`DECISIONS.md`)

## Purpose

This document records every non-trivial architectural decision made during the Go port of `decimal.js`.

Every entry must follow this format:

- **Context:** What problem were we solving?
- **Evidence:** What did we learn from reverse-engineering `decimal.js`? (line numbers, test references, commit history)
- **Alternatives Considered:** What other approaches exist? (with references to `apd`, `shopspring`, `math/big`, `bignumber.js`)
- **Decision:** What did we choose?
- **Rationale:** Why? Based on evidence, not assumption.

---

## Commitments

1. **Zero Unsafe:** 0 `unsafe.Pointer`, 0 `any` / `interface{}` escape hatches across all Go code.
2. **Zero Test Tampering:** Original `decimal.js` test suite in `tests/original/` remains 100% untouched.
3. **No decision without evidence:** Every entry below was earned through reverse engineering, not assumed.

---

## Decisions

*No decisions have been made yet.*

*Decisions will be recorded here as the research phase (`RESEARCH_PROTOCOL.md`, `REPOSITORY_INDEX.md`, `COMPARATIVE_ANALYSIS.md`) produces evidence that justifies them.*

*The following are examples of decisions that WILL need to be made, but only after research:*

- Internal representation (coefficient type, base, exponent type)
- API surface design (how to translate `clone()` to Go)
- Rounding implementation strategy
- Algorithm selection for each advanced operation
- Memory allocation strategy
- Error handling pattern
- String parsing approach

*Each will be documented here with full evidence when the time comes.*
