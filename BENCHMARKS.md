# Benchmarking Guide

This document describes how to run, profile, and interpret benchmarks for the mapper library.

## Running Benchmarks

### Basic Benchmark Run

```bash
go test -bench="." -benchmem ./...
```

> **Note**: The quotes around `.` are required on Windows (PowerShell/cmd). On Unix systems, `go test -bench=. -benchmem ./...` also works.

### Run Specific Benchmark Pattern

```bash
# Run only flat struct benchmarks
go test -bench=Flat -benchmem ./...

# Run only slice benchmarks
go test -bench=Slice -benchmem ./...

# Run only baseline comparisons
go test -bench=Baseline -benchmem ./...

# Run parallel benchmarks
go test -bench=Parallel -benchmem ./...
```

### Extended Benchmark Run (More Iterations)

For more stable results, increase the benchmark time:

```bash
go test -bench="." -benchmem -benchtime=5s ./...
```

### Run Benchmarks Multiple Times for Statistical Analysis

```bash
# Run 10 times and save results for benchstat
go test -bench="." -benchmem -count=10 ./... > bench_results.txt
```

## CPU and Memory Profiling

### Generate CPU Profile

```bash
# Run benchmarks with CPU profiling
go test -bench=BenchmarkMap_Flat -benchmem -cpuprofile=cpu.prof ./...

# Analyze the profile
go tool pprof cpu.prof
```

Once inside pprof interactive mode:
- `top10` - Show top 10 functions by CPU time
- `list Map` - Show source code for Map function
- `web` - Open interactive graph in browser (requires graphviz)

### Generate Memory Profile

```bash
# Run benchmarks with memory profiling
go test -bench=BenchmarkMap_Flat -benchmem -memprofile=mem.prof ./...

# Analyze allocations
go tool pprof -alloc_space mem.prof

# Analyze objects in use
go tool pprof -inuse_space mem.prof
```

Once inside pprof interactive mode:
- `top10` - Show top 10 allocating functions
- `list assignValue` - Show allocations in assignValue function

### Generate Both Profiles

```bash
go test -bench=BenchmarkMap_Nested -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof ./...
```

### Trace Generation (for Scheduler Analysis)

```bash
go test -bench=BenchmarkMap_Parallel -trace=trace.out ./...
go tool trace trace.out
```

## Comparing Results with benchstat

benchstat is the standard tool for comparing benchmark results. Install it:

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

### Compare Two Benchmark Runs

```bash
# Run baseline
go test -bench=. -benchmem -count=10 ./... > old.txt

# Make changes...

# Run again
go test -bench=. -benchmem -count=10 ./... > new.txt

# Compare
benchstat old.txt new.txt
```

### Example benchstat Output

```
name                          old time/op    new time/op    delta
Map_Flat-8                      1.23µs ± 2%    1.15µs ± 1%   -6.50%  (p=0.000 n=10+10)
Map_Nested-8                    2.45µs ± 3%    2.38µs ± 2%   -2.86%  (p=0.001 n=10+10)

name                          old alloc/op   new alloc/op   delta
Map_Flat-8                        432B ± 0%      400B ± 0%   -7.41%  (p=0.000 n=10+10)
Map_Nested-8                    1.02kB ± 0%    0.98kB ± 0%   -3.92%  (p=0.000 n=10+10)

name                          old allocs/op  new allocs/op  delta
Map_Flat-8                        12.0 ± 0%      11.0 ± 0%   -8.33%  (p=0.000 n=10+10)
Map_Nested-8                      28.0 ± 0%      26.0 ± 0%   -7.14%  (p=0.000 n=10+10)
```

## Interpreting Results

### Understanding Benchmark Output

```
BenchmarkMap_Flat-8         500000      2341 ns/op      432 B/op      12 allocs/op
```

| Field | Meaning |
|-------|---------|
| `BenchmarkMap_Flat-8` | Benchmark name, `-8` indicates GOMAXPROCS=8 |
| `500000` | Number of iterations run |
| `2341 ns/op` | Nanoseconds per operation |
| `432 B/op` | Bytes allocated per operation |
| `12 allocs/op` | Number of heap allocations per operation |

### Performance Guidelines

**ns/op (Latency)**
- < 1µs: Excellent for simple mappings
- 1-10µs: Acceptable for typical DTO mapping
- > 10µs: May need investigation for hot paths

**allocs/op (Allocation Count)**
- Lower is better; each allocation adds GC pressure
- Flat struct mapping should be < 20 allocations
- Slice mapping scales linearly with slice length

**B/op (Memory)**
- Should be proportional to the data being copied
- Watch for unexpected growth indicating memory leaks

### Baseline Comparison

Compare reflection-based mapping against manual mapping:

```bash
go test -bench=Baseline -benchmem ./...
go test -bench=Map_Flat -benchmem ./...
```

Expected overhead of reflection-based mapping:
- **10-50x** slower than manual mapping (typical for reflection)
- **2-5x** more allocations than manual mapping

The tradeoff is development time vs runtime performance. For most applications, the reflection overhead is acceptable compared to network I/O or database latency.

### Parallel Benchmark Interpretation

Parallel benchmarks reveal lock contention. Compare:

```bash
# Sequential
go test -bench=Map_Flat -benchmem ./...

# Parallel
go test -bench=Map_Parallel_Flat -benchmem ./...
```

If parallel performance doesn't scale linearly with GOMAXPROCS, investigate:
1. Lock contention in the metadata cache
2. False sharing in CPU caches
3. Memory allocator contention

## CI Integration

### GitHub Actions Example

```yaml
jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Run Benchmarks
        run: go test -bench=. -benchmem -count=5 ./... > benchmark.txt
      - name: Upload Results
        uses: actions/upload-artifact@v4
        with:
          name: benchmark-results
          path: benchmark.txt
```

### Storing Baseline Results

```bash
# After a release, save the baseline
go test -bench=. -benchmem -count=10 ./... > benchmark_baseline.txt
git add benchmark_baseline.txt
git commit -m "chore: update benchmark baseline"
```

## Benchmark Categories

| Benchmark | Purpose |
|-----------|---------|
| `Map_Flat` | Baseline for simple struct mapping |
| `Map_Nested` | Tests recursive struct traversal |
| `Map_DeepNested` | Stress test for deep nesting |
| `Map_Slice_*` | Tests slice allocation and copying |
| `Map_Optional_*` | Tests pointer handling paths |
| `MapWithOptions_*` | Tests configuration overhead |
| `Map_CacheWarm` | Tests cached metadata performance |
| `Baseline_Manual*` | Reference for zero-overhead mapping |
| `Map_Parallel_*` | Tests thread safety and contention |

## Troubleshooting

### Inconsistent Results

If benchmarks show high variance (± > 5%):

1. Close other applications
2. Disable CPU frequency scaling
3. Run on an isolated machine
4. Increase `-benchtime` or `-count`

### Memory Profile Shows Unexpected Allocations

1. Use `go tool pprof -alloc_objects` to count allocations
2. Look for allocations in reflect package (expected)
3. Check for string concatenation in hot paths
4. Verify slice pre-allocation

### CPU Profile Shows Hot Spots

Common hot spots in reflection-based mappers:
- `reflect.Value.Interface()` - boxing/unboxing
- `reflect.Value.Set()` - type checking
- Map/slice iteration with reflection

These are inherent to reflection and difficult to optimize without code generation.
