package iter_test

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
)

// ============================================================================
// MapLazy vs Map (eager) Benchmarks
// ============================================================================

func BenchmarkMapLazy_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		mapped := iter.MapLazy(it, func(v int) int { return v * 2 })
		_ = mapped.Collect()
	}
}

func BenchmarkMapEager_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		mapped := iter.Map(it, func(v int) int { return v * 2 })
		_ = mapped.Count()
	}
}

func BenchmarkMapLazy_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		mapped := iter.MapLazy(it, func(v int) int { return v * 2 })
		_ = mapped.Collect()
	}
}

func BenchmarkMapEager_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		mapped := iter.Map(it, func(v int) int { return v * 2 })
		_ = mapped.Count()
	}
}

// ============================================================================
// TakeLazy vs Take (eager) Benchmarks
// ============================================================================

func BenchmarkTakeLazy_100_from_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		taken := it.TakeLazy(100)
		_ = taken.Collect()
	}
}

func BenchmarkTakeEager_100_from_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		taken := it.Take(100)
		_ = taken.Count()
	}
}

// ============================================================================
// FilterMapLazy vs FilterMap (eager) Benchmarks
// ============================================================================

func BenchmarkFilterMapLazy_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		fm := iter.FilterMapLazy(it, func(v int) shepard.Option[int] {
			if v%2 == 0 {
				return shepard.Some(v * 2)
			}
			return shepard.None[int]()
		})
		_ = fm.Collect()
	}
}

func BenchmarkFilterMapEager_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		fm := iter.FilterMap(it, func(v int) shepard.Option[int] {
			if v%2 == 0 {
				return shepard.Some(v * 2)
			}
			return shepard.None[int]()
		})
		_ = fm.Count()
	}
}

func BenchmarkFilterMapLazy_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		fm := iter.FilterMapLazy(it, func(v int) shepard.Option[int] {
			if v%2 == 0 {
				return shepard.Some(v * 2)
			}
			return shepard.None[int]()
		})
		_ = fm.Collect()
	}
}

func BenchmarkFilterMapEager_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := iter.New(data)
		fm := iter.FilterMap(it, func(v int) shepard.Option[int] {
			if v%2 == 0 {
				return shepard.Some(v * 2)
			}
			return shepard.None[int]()
		})
		_ = fm.Count()
	}
}

// ============================================================================
// Chain Benchmarks
// ============================================================================

func BenchmarkChain_1000_each(b *testing.B) {
	data1 := makeData(1000)
	data2 := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it1 := iter.New(data1)
		it2 := iter.New(data2)
		chained := it1.Chain(it2)
		_ = chained.Collect()
	}
}

func BenchmarkNative_Concat_1000_each(b *testing.B) {
	data1 := makeData(1000)
	data2 := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := make([]int, 0, len(data1)+len(data2))
		result = append(result, data1...)
		result = append(result, data2...)
		_ = len(result)
	}
}
