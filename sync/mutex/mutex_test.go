package mutex_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/marlaone/shepard/sync/mutex"
	"github.com/stretchr/testify/assert"
)

func TestMutex_New(t *testing.T) {
	m := mutex.New[int](42)
	result := m.Lock()
	assert.Equal(t, 42, result.Unwrap())
	m.Unlock()
}

func TestMutex_Lock(t *testing.T) {
	m := mutex.New[uint8](1)
	assert.Equal(t, uint8(1), m.Lock().Unwrap())
	m.Unlock()
}

func TestMutex_Unlock(t *testing.T) {
	m := mutex.New[string]("test")
	_ = m.Lock()
	// Should not panic
	m.Unlock()

	// Should be able to lock again after unlock
	result := m.Lock()
	assert.Equal(t, "test", result.Unwrap())
	m.Unlock()
}

func TestMutex_LockUnlockCycle(t *testing.T) {
	m := mutex.New[int](0)

	for i := 0; i < 100; i++ {
		result := m.Lock()
		assert.Equal(t, 0, result.Unwrap())
		m.Unlock()
	}
}

func TestMutex_MutualExclusion(t *testing.T) {
	m := mutex.New[int](0)
	var counter int32
	var wg sync.WaitGroup

	// Multiple goroutines try to acquire lock
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = m.Lock()
			// Critical section - increment counter
			current := atomic.LoadInt32(&counter)
			atomic.StoreInt32(&counter, current+1)
			// Small delay to increase chance of race if mutex fails
			time.Sleep(time.Millisecond)
			m.Unlock()
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(10), atomic.LoadInt32(&counter))
}

func TestMutex_BlocksWhileLocked(t *testing.T) {
	m := mutex.New[int](42)
	var acquired atomic.Bool
	var wg sync.WaitGroup

	// First goroutine holds the lock
	_ = m.Lock()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// This should block until first goroutine unlocks
		_ = m.Lock()
		acquired.Store(true)
		m.Unlock()
	}()

	// Give the goroutine time to attempt acquiring the lock
	time.Sleep(50 * time.Millisecond)
	assert.False(t, acquired.Load(), "Second goroutine should be blocked")

	// Now unlock and let the second goroutine proceed
	m.Unlock()
	wg.Wait()
	assert.True(t, acquired.Load(), "Second goroutine should have acquired lock")
}

func TestMutex_DifferentTypes(t *testing.T) {
	// Test with struct type
	type Data struct {
		Value int
		Name  string
	}

	m := mutex.New[Data](Data{Value: 10, Name: "test"})
	result := m.Lock()
	data := result.Unwrap()
	assert.Equal(t, 10, data.Value)
	assert.Equal(t, "test", data.Name)
	m.Unlock()

	// Test with slice type
	m2 := mutex.New[[]int]([]int{1, 2, 3})
	result2 := m2.Lock()
	assert.Equal(t, []int{1, 2, 3}, result2.Unwrap())
	m2.Unlock()
}
