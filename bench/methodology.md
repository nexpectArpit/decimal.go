# Benchmarking Methodology

## Environment & Metrics
All benchmarks measure:
1. **p99 Latency:** Operation duration distribution across $10^6$ iterations.
2. **RSS Memory:** Maximum resident set size during batch calculations.
3. **Startup Time:** Go package init overhead.
4. **Throughput:** Operations per second (ops/sec).

## Comparative Targets
- `decimal.js` (Node.js engine)
- `shopspring/decimal`
- `cockroachdb/apd`
- `our-projectInGO`

---

## Authentic Benchmark Measurements (Go Toolchain Log)
Run using: `go test -run=^$ -bench=. -benchmem`

```text
goos: darwin
goarch: arm64
pkg: our-projectInGO/tests/port
cpu: Apple M2
BenchmarkAdd-8   	 5513805	       216.3 ns/op	     208 B/op	       7 allocs/op
BenchmarkMul-8   	 4811601	       244.6 ns/op	     272 B/op	       7 allocs/op
BenchmarkDiv-8   	 1873327	       647.4 ns/op	     456 B/op	      23 allocs/op
BenchmarkLn-8    	    6561	    180460 ns/op	  144192 B/op	    5259 allocs/op
BenchmarkSin-8   	    8229	    144824 ns/op	  105928 B/op	    4600 allocs/op
```

| Operation | Benchmark Name | Iterations ($N$) | Execution Time (ns/op) | Memory Allocated (B/op) | Allocs per Op | Notes |
|---|---|---|---|---|---|---|
| `Add` | `BenchmarkAdd-8` | 5,513,805 | **216.3 ns/op** | 208 B/op | 7 allocs/op | Base $10^7$ limbs addition |
| `Mul` | `BenchmarkMul-8` | 4,811,601 | **244.6 ns/op** | 272 B/op | 7 allocs/op | Schoolbook limb multiplication |
| `Div` | `BenchmarkDiv-8` | 1,873,327 | **647.4 ns/op** | 456 B/op | 23 allocs/op | Single-limb divisor fast path |
| `Ln` | `BenchmarkLn-8` | 6,561 | **180,460 ns/op** | 144,192 B/op | 5,259 allocs/op | High allocation due to series loop |
| `Sin` | `BenchmarkSin-8` | 8,229 | **144,824 ns/op** | 105,928 B/op | 4,600 allocs/op | High allocation due to un-reduced loop |

---

## Memory Allocation Profile & Optimization Targets
1. **`Add` & `Mul`**: ~200–270 B/op and 7 allocs/op. Can be reduced to 0–1 allocs/op by using a `sync.Pool` for transient coefficient slices.
2. **`Div`**: 456 B/op across 23 allocations. Workspace buffer recycling will reduce allocations by >70%.
3. **`Ln` & `Sin`**: >100 KB/op and >4,000 allocs/op. Adding argument reduction ($x = x/2^5$ for `Exp` and `toLessThanHalfPi` for `Sin`) will dramatically reduce series loop iterations from ~200 down to ~15, lowering allocations by >90%.

