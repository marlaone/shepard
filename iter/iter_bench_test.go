package iter_test

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
)

// Helper to create test data
func makeData(n int) []int {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = i
	}
	return data
}

// ============================================================================
// Foreach Benchmarks
// ============================================================================

func BenchmarkIter_Foreach_100(b *testing.B) {
	data := makeData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		it.Foreach(func(idx int, val int) {
			_ = val
		})
	}
}

func BenchmarkNative_Range_100(b *testing.B) {
	data := makeData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range data {
			_ = v
		}
	}
}

func BenchmarkIter_Foreach_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		it.Foreach(func(idx int, val int) {
			_ = val
		})
	}
}

func BenchmarkNative_Range_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range data {
			_ = v
		}
	}
}

func BenchmarkIter_Foreach_100000(b *testing.B) {
	data := makeData(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		it.Foreach(func(idx int, val int) {
			_ = val
		})
	}
}

func BenchmarkNative_Range_100000(b *testing.B) {
	data := makeData(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range data {
			_ = v
		}
	}
}

// ============================================================================
// Filter Benchmarks (keep even numbers)
// ============================================================================

func BenchmarkIter_Filter_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		filtered := it.Filter(func(v *int) bool { return *v%2 == 0 })
		_ = filtered.Count()
	}
}

func BenchmarkNative_Filter_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []int
		for _, v := range data {
			if v%2 == 0 {
				result = append(result, v)
			}
		}
		_ = len(result)
	}
}

func BenchmarkIter_Filter_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		filtered := it.Filter(func(v *int) bool { return *v%2 == 0 })
		_ = filtered.Count()
	}
}

func BenchmarkNative_Filter_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []int
		for _, v := range data {
			if v%2 == 0 {
				result = append(result, v)
			}
		}
		_ = len(result)
	}
}

// ============================================================================
// Map Benchmarks (double each value)
// ============================================================================

func BenchmarkIter_Map_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		mapped := iter.Map(it, func(v int) int { return v * 2 })
		_ = mapped.Count()
	}
}

func BenchmarkNative_Map_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, len(data))
		for j, v := range data {
			result[j] = v * 2
		}
		_ = len(result)
	}
}

func BenchmarkIter_Map_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		mapped := iter.Map(it, func(v int) int { return v * 2 })
		_ = mapped.Count()
	}
}

func BenchmarkNative_Map_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, len(data))
		for j, v := range data {
			result[j] = v * 2
		}
		_ = len(result)
	}
}

// ============================================================================
// FilterMap Benchmarks (filter evens, then double)
// ============================================================================

func BenchmarkIter_FilterMap_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		result := iter.FilterMap(it, func(v int) shepard.Option[int] {
			if v%2 == 0 {
				return shepard.Some(v * 2)
			}
			return shepard.None[int]()
		})
		_ = result.Count()
	}
}

func BenchmarkNative_FilterMap_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []int
		for _, v := range data {
			if v%2 == 0 {
				result = append(result, v*2)
			}
		}
		_ = len(result)
	}
}

func BenchmarkIter_FilterMap_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		result := iter.FilterMap(it, func(v int) shepard.Option[int] {
			if v%2 == 0 {
				return shepard.Some(v * 2)
			}
			return shepard.None[int]()
		})
		_ = result.Count()
	}
}

func BenchmarkNative_FilterMap_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []int
		for _, v := range data {
			if v%2 == 0 {
				result = append(result, v*2)
			}
		}
		_ = len(result)
	}
}

// ============================================================================
// Find Benchmarks (find first > 500)
// ============================================================================

func BenchmarkIter_Find_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		_ = it.Find(func(v *int) bool { return *v > 500 })
	}
}

func BenchmarkNative_Find_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var found *int
		for j := range data {
			if data[j] > 500 {
				found = &data[j]
				break
			}
		}
		_ = found
	}
}

func BenchmarkIter_Find_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		_ = it.Find(func(v *int) bool { return *v > 5000 })
	}
}

func BenchmarkNative_Find_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var found *int
		for j := range data {
			if data[j] > 5000 {
				found = &data[j]
				break
			}
		}
		_ = found
	}
}

// ============================================================================
// Take Benchmarks
// ============================================================================

func BenchmarkIter_Take_100_from_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		taken := it.Take(100)
		_ = taken.Count()
	}
}

func BenchmarkNative_Take_100_from_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := data[:100]
		_ = len(result)
	}
}

func BenchmarkIter_Take_1000_from_100000(b *testing.B) {
	data := makeData(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		taken := it.Take(1000)
		_ = taken.Count()
	}
}

func BenchmarkNative_Take_1000_from_100000(b *testing.B) {
	data := makeData(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := data[:1000]
		_ = len(result)
	}
}

// ============================================================================
// Chained Operations: Filter then Map
// ============================================================================

func BenchmarkIter_Chain_FilterThenMap_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		filtered := it.Filter(func(v *int) bool { return *v%2 == 0 })
		mapped := iter.Map(filtered, func(v int) int { return v * 2 })
		_ = mapped.Count()
	}
}

func BenchmarkNative_Chain_FilterThenMap_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []int
		for _, v := range data {
			if v%2 == 0 {
				result = append(result, v*2)
			}
		}
		_ = len(result)
	}
}

// ============================================================================
// Count Benchmarks
// ============================================================================

func BenchmarkIter_Count_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		_ = it.Count()
	}
}

func BenchmarkNative_Len_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = len(data)
	}
}

// ============================================================================
// Next (manual iteration) Benchmarks
// ============================================================================

func BenchmarkIter_Next_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		for {
			opt := it.Next()
			if opt.IsNone() {
				break
			}
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNative_Index_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(data); j++ {
			_ = data[j]
		}
	}
}

func BenchmarkIter_Next_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		for {
			opt := it.Next()
			if opt.IsNone() {
				break
			}
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNative_Index_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(data); j++ {
			_ = data[j]
		}
	}
}
