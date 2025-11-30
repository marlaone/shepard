# HTTP Profiling Results

Generated: 2025-11-30

## Benchmark Baseline

Key metrics from baseline run (3 iterations, -benchmem):

| Benchmark | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| Router_Match_Static_5Routes | ~805 | 1000 | 19 |
| Router_Match_Static_50Routes | ~4070 | 4680 | 88 |
| Router_Match_Static_200Routes | ~13950 | 16680 | 313 |
| Router_Match_WithParams_1Param | ~820 | 1256 | 20 |
| Router_Match_WithParams_3Params | ~1243 | 1624 | 29 |
| FullCycle_SimpleGet | ~807 | 1184 | 20 |
| FullCycle_GetWithParams | ~1243 | 1760 | 27 |
| FullCycle_PostJSON_Small | ~3558 | 2793 | 43 |
| FullCycle_PostJSON_Large | ~202903 | 113213 | 853 |

## CPU Profile Analysis

### Router Benchmarks

Top functions by cumulative CPU time:
1. `Router.Handle` - 17.28%
2. `Router.findHandler` - 13.27%
3. `RoutePattern.Match` - 7.83%
4. `strings.genSplit` - 4.67%
5. `mallocgc` - 8.79%

### Full Cycle Benchmarks

Top functions (excluding runtime):
1. `encoding/json.(*decodeState).object` - 2.88%
2. `encoding/json.checkValid` - 1.39%
3. `encoding/json.(*decodeState).value` - 2.90%

**Insight:** JSON parsing dominates CPU in POST benchmarks. Router matching is the primary overhead in GET requests.

## Memory Profile Analysis

### Router Allocation Sources (by object count)

| Function | % of Allocations |
|----------|------------------|
| Router.findHandler | 24.16% |
| strings.genSplit | 21.65% |
| hashmap.New[string,string] | 12.76% |
| HashMap.Insert (params) | 10.06% |
| HashMap.Insert (headers) | 9.91% |
| NewBodyWriter | 3.38% |
| NewHeaders | 3.36% |

### Full Cycle Allocation Sources

| Function | % of Allocations |
|----------|------------------|
| HashMap.Insert (headers) | 14.46% |
| HashMap.New (headers) | 6.60% |
| HashMap.New (params) | 6.53% |
| Router.findHandler | 6.37% |
| NewHeaders | 6.28% |
| slice.Init[string] | 6.12% |
| strings.Builder.grow | 5.95% |
| HashMap.Insert (params) | 4.13% |

## Identified Hotspots

### 1. RoutePattern.Match (CRITICAL)
- **File:** `http/route.go`
- **Issue:** Uses `strings.Split` on every match, allocating string slices
- **Allocations:** ~40% of router allocations
- **Fix:** Index-based parsing without allocation

### 2. Path Parameters HashMap (HIGH)
- **File:** `http/router.go:152-157`, `http/context.go`
- **Issue:** New HashMap created per route match
- **Allocations:** ~13% of total
- **Fix:** sync.Pool or fixed-size array

### 3. Response Headers (HIGH)
- **File:** `http/headers.go`
- **Issue:** New HashMap per response
- **Allocations:** ~10% of total
- **Fix:** sync.Pool for Headers

### 4. BodyWriter (MEDIUM)
- **File:** `http/body.go`
- **Issue:** New allocation per response
- **Allocations:** ~3.4% of total
- **Fix:** sync.Pool

### 5. Header Canonicalization (MEDIUM)
- **File:** `http/headers.go`
- **Issue:** strings.Builder.grow for each header
- **Allocations:** ~5-6% of total
- **Fix:** Cache common canonical names

## Optimization Recommendations

### Quick Wins (Low Risk, High Impact)

1. **Add sync.Pool for BodyWriter**
   - File: `http/body.go`
   - Expected: ~3% allocation reduction

2. **Add sync.Pool for Headers**
   - File: `http/headers.go`
   - Expected: ~10% allocation reduction

3. **Cache canonical header names**
   - File: `http/headers.go`
   - Expected: ~5% allocation reduction

4. **Pre-allocate path params capacity**
   - File: `http/context.go`
   - Expected: Reduced grow operations

### Medium Effort (Good ROI)

5. **Zero-allocation path matching**
   - File: `http/route.go`
   - Expected: ~40% allocation reduction in route matching
   - Use index-based parsing instead of strings.Split

6. **Reuse path segment slices with pool**
   - File: `http/route.go`
   - Alternative to #5 if full rewrite is too risky

### Major Refactors (Consider Later)

7. **Trie-based router**
   - File: `http/router.go`
   - Impact: O(path_length) instead of O(num_routes)
   - Would eliminate 200-route scaling problem

8. **Streaming file responses**
   - File: `http/into_response.go`
   - Impact: Avoid buffering entire files in memory

## Target Metrics After Optimization

Based on analysis, realistic targets:

| Metric | Current | Target | Improvement |
|--------|---------|--------|-------------|
| Router_200Routes allocs | 313 | ~150 | -50% |
| Router_200Routes ns/op | 13950 | ~8000 | -40% |
| FullCycle_SimpleGet allocs | 20 | ~12 | -40% |
| Path params extraction | 3 allocs | 1 alloc | -66% |

## Profiling Commands Reference

```bash
# CPU profiling
go test -bench=BenchmarkRouter -cpuprofile=cpu.prof -benchtime=5s ./http/benchmarks/

# Memory profiling (allocation count)
go tool pprof -alloc_objects -top mem.prof

# Memory profiling (allocation size)
go tool pprof -alloc_space -top mem.prof

# Interactive analysis
go tool pprof -http=:8080 cpu.prof

# Compare before/after
benchstat bench_before.txt bench_after.txt
```
