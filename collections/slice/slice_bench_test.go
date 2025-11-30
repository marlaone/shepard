package slice_test

import (
	"testing"

	"github.com/marlaone/shepard/collections/slice"
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
// Push / Append Benchmarks
// ============================================================================

func BenchmarkSlice_Push_100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := slice.New[int]()
		for j := 0; j < 100; j++ {
			s.Push(j)
		}
	}
}

func BenchmarkNative_Append_100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []int
		for j := 0; j < 100; j++ {
			s = append(s, j)
		}
	}
}

func BenchmarkSlice_Push_10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := slice.New[int]()
		for j := 0; j < 10000; j++ {
			s.Push(j)
		}
	}
}

func BenchmarkNative_Append_10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []int
		for j := 0; j < 10000; j++ {
			s = append(s, j)
		}
	}
}

func BenchmarkSlice_Push_100000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := slice.New[int]()
		for j := 0; j < 100000; j++ {
			s.Push(j)
		}
	}
}

func BenchmarkNative_Append_100000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var s []int
		for j := 0; j < 100000; j++ {
			s = append(s, j)
		}
	}
}

// ============================================================================
// Pop Benchmarks
// ============================================================================

func BenchmarkSlice_Pop_100(b *testing.B) {
	data := makeData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		for s.Len() > 0 {
			_ = s.Pop()
		}
	}
}

func BenchmarkNative_Pop_100(b *testing.B) {
	data := makeData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		for len(s) > 0 {
			_ = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}
}

func BenchmarkSlice_Pop_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		for s.Len() > 0 {
			_ = s.Pop()
		}
	}
}

func BenchmarkNative_Pop_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		for len(s) > 0 {
			_ = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}
}

// ============================================================================
// Get / Index Access Benchmarks
// ============================================================================

func BenchmarkSlice_Get_10000(b *testing.B) {
	s := slice.Init(makeData(10000)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			_ = s.Get(j)
		}
	}
}

func BenchmarkNative_Index_10000(b *testing.B) {
	s := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10000; j++ {
			_ = s[j]
		}
	}
}

// ============================================================================
// Iteration Benchmarks
// ============================================================================

func BenchmarkSlice_Iter_100(b *testing.B) {
	s := slice.Init(makeData(100)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := s.Iter()
		for {
			opt := iter.Next()
			if opt.IsNone() {
				break
			}
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNative_Range_100(b *testing.B) {
	s := makeData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range s {
			_ = v
		}
	}
}

func BenchmarkSlice_Iter_10000(b *testing.B) {
	s := slice.Init(makeData(10000)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := s.Iter()
		for {
			opt := iter.Next()
			if opt.IsNone() {
				break
			}
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNative_Range_10000(b *testing.B) {
	s := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range s {
			_ = v
		}
	}
}

func BenchmarkSlice_Iter_100000(b *testing.B) {
	s := slice.Init(makeData(100000)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iter := s.Iter()
		for {
			opt := iter.Next()
			if opt.IsNone() {
				break
			}
			_ = opt.Unwrap()
		}
	}
}

func BenchmarkNative_Range_100000(b *testing.B) {
	s := makeData(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range s {
			_ = v
		}
	}
}

// ============================================================================
// Insert Benchmarks (Middle)
// ============================================================================

func BenchmarkSlice_Insert_Middle_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		s.Insert(500, 999)
	}
}

func BenchmarkNative_Insert_Middle_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		idx := 500
		s = append(s[:idx+1], s[idx:]...)
		s[idx] = 999
	}
}

func BenchmarkSlice_Insert_Middle_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		s.Insert(5000, 999)
	}
}

func BenchmarkNative_Insert_Middle_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		idx := 5000
		s = append(s[:idx+1], s[idx:]...)
		s[idx] = 999
	}
}

// ============================================================================
// Remove Benchmarks (Middle)
// ============================================================================

func BenchmarkSlice_Remove_Middle_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		_ = s.Remove(500)
	}
}

func BenchmarkNative_Remove_Middle_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		idx := 500
		_ = s[idx]
		s = append(s[:idx], s[idx+1:]...)
	}
}

func BenchmarkSlice_Remove_Middle_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		_ = s.Remove(5000)
	}
}

func BenchmarkNative_Remove_Middle_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		idx := 5000
		_ = s[idx]
		s = append(s[:idx], s[idx+1:]...)
	}
}

// ============================================================================
// Clone Benchmarks
// ============================================================================

func BenchmarkSlice_Clone_1000(b *testing.B) {
	s := slice.Init(makeData(1000)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Clone()
	}
}

func BenchmarkNative_Clone_1000(b *testing.B) {
	s := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := make([]int, len(s))
		copy(clone, s)
		_ = clone
	}
}

func BenchmarkSlice_Clone_100000(b *testing.B) {
	s := slice.Init(makeData(100000)...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Clone()
	}
}

func BenchmarkNative_Clone_100000(b *testing.B) {
	s := makeData(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := make([]int, len(s))
		copy(clone, s)
		_ = clone
	}
}

// ============================================================================
// Retain / Filter Benchmarks (keep even numbers)
// ============================================================================

func BenchmarkSlice_Retain_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		s.Retain(func(x *int) bool { return *x%2 == 0 })
	}
}

func BenchmarkNative_Filter_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		result := s[:0]
		for _, v := range s {
			if v%2 == 0 {
				result = append(result, v)
			}
		}
		_ = result
	}
}

func BenchmarkSlice_Retain_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		s.Retain(func(x *int) bool { return *x%2 == 0 })
	}
}

func BenchmarkNative_Filter_10000(b *testing.B) {
	data := makeData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		result := s[:0]
		for _, v := range s {
			if v%2 == 0 {
				result = append(result, v)
			}
		}
		_ = result
	}
}

// ============================================================================
// Reverse Benchmarks
// ============================================================================

func BenchmarkSlice_Reverse_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		s.Reverse()
	}
}

func BenchmarkNative_Reverse_1000(b *testing.B) {
	data := makeData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		for left, right := 0, len(s)-1; left < right; left, right = left+1, right-1 {
			s[left], s[right] = s[right], s[left]
		}
	}
}

func BenchmarkSlice_Reverse_100000(b *testing.B) {
	data := makeData(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := slice.Init(data...)
		s.Reverse()
	}
}

func BenchmarkNative_Reverse_100000(b *testing.B) {
	data := makeData(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, len(data))
		copy(s, data)
		for left, right := 0, len(s)-1; left < right; left, right = left+1, right-1 {
			s[left], s[right] = s[right], s[left]
		}
	}
}
