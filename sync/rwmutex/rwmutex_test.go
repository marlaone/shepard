package rwmutex_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/marlaone/shepard/sync/rwmutex"
	"github.com/stretchr/testify/assert"
)

func TestRWMutex_New(t *testing.T) {
	m := rwmutex.New[int](42)
	result := m.Lock()
	assert.Equal(t, 42, *result.Unwrap())
	m.Unlock()
}

func TestRWMutex_Lock(t *testing.T) {
	m := rwmutex.New[string]("hello")

	// Acquire write lock
	guard := m.Lock()
	value := guard.Unwrap()
	assert.Equal(t, "hello", *value)

	// Modify through pointer
	*value = "world"

	m.Unlock()

	// Verify change persisted
	guard2 := m.Lock()
	assert.Equal(t, "world", *guard2.Unwrap())
	m.Unlock()
}

func TestRWMutex_RLock(t *testing.T) {
	m := rwmutex.New[int](100)

	// Acquire read lock
	guard := m.RLock()
	value := guard.Unwrap()
	assert.Equal(t, 100, *value)
	m.RUnlock()
}

func TestRWMutex_ConcurrentReads(t *testing.T) {
	m := rwmutex.New[int](0)
	var wg sync.WaitGroup
	var readCount atomic.Int32

	// Start 10 concurrent readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			guard := m.RLock()
			readCount.Add(1)
			_ = guard.Unwrap()
			// Hold the read lock briefly to allow overlap
			time.Sleep(10 * time.Millisecond)
			m.RUnlock()
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(10), readCount.Load())
}

func TestRWMutex_WriteLockBlocksReads(t *testing.T) {
	m := rwmutex.New[int](0)
	var sequence []string
	var mu sync.Mutex

	appendSeq := func(s string) {
		mu.Lock()
		sequence = append(sequence, s)
		mu.Unlock()
	}

	// Acquire write lock first
	m.Lock()
	appendSeq("write_start")

	var wg sync.WaitGroup
	wg.Add(1)

	// Start a reader that should block
	go func() {
		defer wg.Done()
		appendSeq("read_waiting")
		m.RLock()
		appendSeq("read_acquired")
		m.RUnlock()
	}()

	// Give the reader time to start waiting
	time.Sleep(50 * time.Millisecond)
	appendSeq("write_end")
	m.Unlock()

	wg.Wait()

	// Reader should have acquired after write released
	assert.Equal(t, "write_start", sequence[0])
	assert.Equal(t, "read_waiting", sequence[1])
	assert.Equal(t, "write_end", sequence[2])
	assert.Equal(t, "read_acquired", sequence[3])
}

func TestRWMutex_ReadLocksBlockWrite(t *testing.T) {
	m := rwmutex.New[int](0)
	var sequence []string
	var mu sync.Mutex

	appendSeq := func(s string) {
		mu.Lock()
		sequence = append(sequence, s)
		mu.Unlock()
	}

	// Acquire read lock first
	m.RLock()
	appendSeq("read_start")

	var wg sync.WaitGroup
	wg.Add(1)

	// Start a writer that should block
	go func() {
		defer wg.Done()
		appendSeq("write_waiting")
		m.Lock()
		appendSeq("write_acquired")
		m.Unlock()
	}()

	// Give the writer time to start waiting
	time.Sleep(50 * time.Millisecond)
	appendSeq("read_end")
	m.RUnlock()

	wg.Wait()

	// Writer should have acquired after read released
	assert.Equal(t, "read_start", sequence[0])
	assert.Equal(t, "write_waiting", sequence[1])
	assert.Equal(t, "read_end", sequence[2])
	assert.Equal(t, "write_acquired", sequence[3])
}

func TestRWMutex_LockUnlockCycle(t *testing.T) {
	m := rwmutex.New[int](0)

	// Multiple lock/unlock cycles
	for i := 0; i < 5; i++ {
		guard := m.Lock()
		*guard.Unwrap() = i
		m.Unlock()
	}

	// Verify final value
	guard := m.RLock()
	assert.Equal(t, 4, *guard.Unwrap())
	m.RUnlock()
}

func TestRWMutex_DataIntegrity(t *testing.T) {
	type Data struct {
		Value int
		Name  string
	}

	m := rwmutex.New[Data](Data{Value: 0, Name: "initial"})
	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			guard := m.Lock()
			data := guard.Unwrap()
			data.Value = val
			data.Name = "updated"
			m.Unlock()
		}(i)
	}

	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			guard := m.RLock()
			data := guard.Unwrap()
			// Just access the data to ensure no race
			_ = data.Value
			_ = data.Name
			m.RUnlock()
		}()
	}

	wg.Wait()

	// Final read to verify data is consistent
	guard := m.RLock()
	data := guard.Unwrap()
	assert.Equal(t, "updated", data.Name)
	m.RUnlock()
}
