# Research Protocol

## Absolute Rule

**No Go code may be written until every question in this document is answered with evidence.**

Evidence means:
- Line numbers in `decimal.js`
- Specific commit hashes
- Test file references
- Algorithm citations
- Concrete data from the source code

Not assumptions. Not "probably." Not "typically."

---

## Phase 0: Repository Indexing

Before any analysis, create a complete index of `decimal.js`.

Every function. Every helper. Every constant. Every internal dependency.

This index lives in `REPOSITORY_INDEX.md`.

Until that index is complete, no analysis begins.

---

## Phase 1: Core Representation Questions

These must be answered by reading `decimal.js` source code directly.

1. How is a Decimal value represented internally?
   - What fields exist on a Decimal instance?
   - What does each field store?
   - What are the valid ranges for each field?
   - How are special values (NaN, ±Infinity, -0) represented?

2. What is BASE and why?
   - What is the exact value of BASE?
   - Where is it defined?
   - Why was this value chosen over alternatives?
   - Where does BASE appear in multiplication?
   - Where does BASE appear in division?
   - Where does BASE appear in string conversion?

3. How are coefficients structured?
   - What does each element of the coefficient array represent?
   - How many digits per element?
   - What is the normalization invariant?
   - When are leading zeros allowed?
   - When are trailing zeros preserved?

4. How does the exponent work?
   - What does the exponent represent relative to the coefficient array?
   - How does it relate to the decimal point position?
   - What are minE and maxE?
   - What happens at the boundaries?

---

## Phase 2: Constructor & Parsing Questions

5. How does the constructor parse input?
   - What input types are accepted?
   - What regex patterns are used for validation?
   - How are hex, octal, binary strings parsed?
   - How are exponential notation strings parsed?
   - What errors are thrown for invalid input?
   - How is sign determined?
   - How is the coefficient array built from a string?

6. How does `clone()` work?
   - What does it return?
   - What state is isolated?
   - What state is shared?
   - How do cloned constructors interact?
   - What is the relationship between a clone and the parent?

7. How does `config()` / `set()` work?
   - What properties can be configured?
   - What are the valid ranges?
   - What validation occurs?
   - Where is the configuration stored?
   - Who reads the configuration?

---

## Phase 3: Arithmetic Questions

For each operation (add, subtract, multiply, divide, modulo):

8. What is the exact algorithm?
   - Not "schoolbook" or "Karatsuba." The exact code path.
   - How are operands aligned?
   - How is carry propagated?
   - How are signs handled?
   - What is the intermediate precision?
   - When does rounding occur?
   - What guard digits are used?

9. What are the special cases?
   - What happens with NaN?
   - What happens with ±Infinity?
   - What happens with zero?
   - What happens with equal magnitude subtraction?

---

## Phase 4: Rounding Questions

10. How does rounding work internally?
    - Where is the rounding function defined?
    - What are its exact parameters?
    - How does it determine the rounding digit?
    - How does it determine the tie-breaking direction?
    - How does it truncate the coefficient array?
    - How does it handle carry from rounding (e.g., 9.9999 rounds to 10)?

11. How do the 9 rounding modes differ?
    - Exact code path for each mode.
    - Exact tie-breaking logic for each mode.

---

## Phase 5: Advanced Operation Questions

For each: sqrt, cbrt, ln, exp, pow, sin, cos, tan, asin, acos, atan, sinh, cosh, tanh, asinh, acosh, atanh:

12. What algorithm is used?
    - Not a guess. The exact algorithm from the source code.
    - What convergence criterion?
    - How is precision managed during iteration?
    - What guard digits?
    - What initial approximation?
    - What series expansion (if any)?
    - What identity transformations?

---

## Phase 6: Formatting & Conversion Questions

13. How does `toString()` work?
    - When does it use exponential notation?
    - How is the string built from coefficients?
    - What role do `toExpNeg` and `toExpPos` play?

14. How does `toFixed()` work?
15. How does `toExponential()` work?
16. How does `toPrecision()` work?
17. How does `toNumber()` work?
18. How does `toFraction()` work?

---

## Phase 7: Dependency & Call Graph Questions

19. What is the complete internal call graph?
    - Which public methods call which internal helpers?
    - Which internal helpers call other internal helpers?
    - What are the leaf functions (no further internal calls)?

20. What are the hidden invariants?
    - What assumptions does the code make that are never checked?
    - What preconditions must hold for internal functions?
    - What postconditions are guaranteed?

---

## Phase 8: Test Coverage Questions

21. What tests exist?
    - Complete list of test files.
    - What does each test file cover?
    - What edge cases are tested?
    - What is NOT tested?

---

## Phase 9: Historical Context Questions

22. What do the commits reveal?
    - Major algorithm changes?
    - Bug fixes that reveal edge cases?
    - Performance changes?

23. What do the issues reveal?
    - Precision bugs reported?
    - Rounding bugs reported?
    - Parsing edge cases?

---

## Verification

Before moving to implementation:

- [ ] Every question above has a concrete answer with source references
- [ ] `REPOSITORY_INDEX.md` is complete
- [ ] `COMPARATIVE_ANALYSIS.md` is complete
- [ ] No assumptions remain — only evidence

Only then does Go code begin.
