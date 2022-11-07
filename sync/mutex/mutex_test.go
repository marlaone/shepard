package mutex_test

import (
	"github.com/marlaone/shepard/sync/mutex"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestMutex_Lock(t *testing.T) {
	m := mutex.New[uint8](1)
	(func() {
		assert.Equal(t, uint8(1), m.Lock().Unwrap())
	})()
	runtime.GC()
	assert.Equal(t, uint8(1), m.Lock().Unwrap())
	runtime.GC()
}
