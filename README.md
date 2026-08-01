# Go Port of decimal.js

Production-quality Go implementation of [decimal.js](https://github.com/MikeMcl/decimal.js) for the **Port Mortem Code Resurrection Hackathon**.

**Track F:** Runtime Modernization (JavaScript → Go)

---

## What This Is

A behaviorally equivalent Go port of `decimal.js` — same edge cases, same precision, same rounding, same test suite, same outputs.

Not a new decimal library. A faithful reconstruction.

---

## Quickstart

```bash
make build
make test
make fuzz
make bench
```

Or via Docker:
```bash
docker build -t decimal-go .
docker run --rm decimal-go
```

---

## Project Documents

| Document | Purpose |
| :--- | :--- |
| `AGENDA.md` | Project charter, phases, scoring criteria |
| `RESEARCH_PROTOCOL.md` | Research questions that must be answered before coding |
| `REPOSITORY_INDEX.md` | Complete index of every function in decimal.js |
| `COMPARATIVE_ANALYSIS.md` | How each reference repo solves the same problems |
| `AUTHOR_PHILOSOPHY.md` | Understanding MikeMcl's engineering thinking |
| `DECISIONS.md` | Architectural decisions with evidence and rationale |

---

## Deliverables

| # | Deliverable | Location |
| :--- | :--- | :--- |
| 01 | Public repo | `our-projectInGO/` |
| 02 | One-step build | `Makefile` / `Dockerfile` |
| 03 | Original tests (untouched) | `tests/original/` |
| 04 | Differential fuzz harness | `fuzz/` |
| 05 | Engineering decisions log | `DECISIONS.md` |
| 06 | Benchmark report | `bench/` |
