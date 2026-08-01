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
