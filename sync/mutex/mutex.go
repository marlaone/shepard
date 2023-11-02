package mutex

import (
	"sync"

	"github.com/marlaone/shepard"
)

type LockResult[T any] shepard.Result[Guard[T], error]

func (l *LockResult[T]) Unwrap() T {
	res := shepard.Result[Guard[T], error](*l)
	return res.Unwrap().t
}

type Guard[T any] struct {
	t T
}

type Mutex[T any] struct {
	mutex *sync.Mutex
	t     T
}

func New[T any](t T) Mutex[T] {
	m := Mutex[T]{
		mutex: &sync.Mutex{},
		t:     t,
	}
	return m
}

func (m *Mutex[T]) Lock() *LockResult[T] {
	m.mutex.Lock()
	guard := Guard[T]{t: m.t}
	res := LockResult[T](shepard.Ok[Guard[T], error](guard))
	return &res
}

func (m *Mutex[T]) Unlock() {
	m.mutex.Unlock()
}
