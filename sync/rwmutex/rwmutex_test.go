package rwmutex_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/marlaone/shepard/sync/mutex"
	"github.com/stretchr/testify/assert"
)

func TestMutex_Lock(t *testing.T) {
	m := mutex.New[uint8](1)
	assert.Equal(t, uint8(1), m.Lock().Unwrap())
	runtime.GC()
	time.Sleep(time.Second)
}
