package iter_test

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/iter"
)

func TestMapLazy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	it := iter.New(data)

	mapped := iter.MapLazy(it, func(x int) int { return x * 2 })

	// Test Count
	if mapped.Count() != 5 {
		t.Errorf("expected Count() = 5, got %d", mapped.Count())
	}

	// Test Collect - count items collected
	result := mapped.Collect()
	expected := []int{2, 4, 6, 8, 10}
	count := 0
	result.Foreach(func(i int, v int) {
		count++
		if i < len(expected) && v != expected[i] {
			t.Errorf("expected result[%d] = %d, got %d", i, expected[i], v)
		}
	})
	if count != len(expected) {
		t.Errorf("expected len = %d, got %d", len(expected), count)
	}
}

func TestMapLazy_Next(t *testing.T) {
	data := []int{1, 2, 3}
	it := iter.New(data)

	mapped := iter.MapLazy(it, func(x int) int { return x * 10 })

	opt := mapped.Next()
	if opt.IsNone() || opt.Unwrap() != 10 {
		t.Errorf("expected 10, got %v", opt)
	}

	opt = mapped.Next()
	if opt.IsNone() || opt.Unwrap() != 20 {
		t.Errorf("expected 20, got %v", opt)
	}

	opt = mapped.Next()
	if opt.IsNone() || opt.Unwrap() != 30 {
		t.Errorf("expected 30, got %v", opt)
	}

	opt = mapped.Next()
	if opt.IsSome() {
		t.Errorf("expected None, got %v", opt)
	}
}

func TestTakeLazy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	it := iter.New(data)

	taken := it.TakeLazy(3)

	// Test Count
	if taken.Count() != 3 {
		t.Errorf("expected Count() = 3, got %d", taken.Count())
	}

	// Test Collect
	result := taken.Collect()
	expected := []int{1, 2, 3}
	count := 0
	result.Foreach(func(i int, v int) {
		count++
		if i < len(expected) && v != expected[i] {
			t.Errorf("expected result[%d] = %d, got %d", i, expected[i], v)
		}
	})
	if count != len(expected) {
		t.Errorf("expected len = %d, got %d", len(expected), count)
	}
}

func TestTakeLazy_MoreThanAvailable(t *testing.T) {
	data := []int{1, 2, 3}
	it := iter.New(data)

	taken := it.TakeLazy(10) // more than available

	if taken.Count() != 3 {
		t.Errorf("expected Count() = 3, got %d", taken.Count())
	}

	result := taken.Collect()
	count := 0
	result.Foreach(func(_ int, _ int) {
		count++
	})
	if count != 3 {
		t.Errorf("expected len = 3, got %d", count)
	}
}

func TestChain(t *testing.T) {
	data1 := []int{1, 2, 3}
	data2 := []int{4, 5, 6}
	it1 := iter.New(data1)
	it2 := iter.New(data2)

	chained := it1.Chain(it2)

	// Test Count
	if chained.Count() != 6 {
		t.Errorf("expected Count() = 6, got %d", chained.Count())
	}

	// Test Collect
	result := chained.Collect()
	expected := []int{1, 2, 3, 4, 5, 6}
	count := 0
	result.Foreach(func(i int, v int) {
		count++
		if i < len(expected) && v != expected[i] {
			t.Errorf("expected result[%d] = %d, got %d", i, expected[i], v)
		}
	})
	if count != len(expected) {
		t.Errorf("expected len = %d, got %d", len(expected), count)
	}
}

func TestChain_Next(t *testing.T) {
	data1 := []int{1, 2}
	data2 := []int{3, 4}
	it1 := iter.New(data1)
	it2 := iter.New(data2)

	chained := it1.Chain(it2)

	for i, expected := range []int{1, 2, 3, 4} {
		opt := chained.Next()
		if opt.IsNone() {
			t.Errorf("expected %d at index %d, got None", expected, i)
		} else if opt.Unwrap() != expected {
			t.Errorf("expected %d at index %d, got %d", expected, i, opt.Unwrap())
		}
	}

	opt := chained.Next()
	if opt.IsSome() {
		t.Errorf("expected None, got %v", opt)
	}
}

func TestFilterMapLazy(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	it := iter.New(data)

	// Keep only even numbers and double them
	fm := iter.FilterMapLazy(it, func(x int) shepard.Option[int] {
		if x%2 == 0 {
			return shepard.Some(x * 2)
		}
		return shepard.None[int]()
	})

	// Test Count
	if fm.Count() != 3 {
		t.Errorf("expected Count() = 3, got %d", fm.Count())
	}

	// Test Collect
	result := fm.Collect()
	expected := []int{4, 8, 12}
	count := 0
	result.Foreach(func(i int, v int) {
		count++
		if i < len(expected) && v != expected[i] {
			t.Errorf("expected result[%d] = %d, got %d", i, expected[i], v)
		}
	})
	if count != len(expected) {
		t.Errorf("expected len = %d, got %d", len(expected), count)
	}
}

func TestFilterMapLazy_Next(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	it := iter.New(data)

	fm := iter.FilterMapLazy(it, func(x int) shepard.Option[string] {
		if x > 3 {
			return shepard.Some("big")
		}
		return shepard.None[string]()
	})

	opt := fm.Next()
	if opt.IsNone() || opt.Unwrap() != "big" {
		t.Errorf("expected 'big', got %v", opt)
	}

	opt = fm.Next()
	if opt.IsNone() || opt.Unwrap() != "big" {
		t.Errorf("expected 'big', got %v", opt)
	}

	opt = fm.Next()
	if opt.IsSome() {
		t.Errorf("expected None, got %v", opt)
	}
}
