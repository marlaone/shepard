package rwmutex

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

type RWMutex[T any] struct {
	mutex *sync.RWMutex
	t     T
}

func New[T any](t T) RWMutex[T] {
	m := RWMutex[T]{
		mutex: &sync.RWMutex{},
		t:     t,
	}
	return m
}

func (m *RWMutex[T]) RLock() *LockResult[*T] {
	m.mutex.Lock()
	guard := Guard[*T]{t: &m.t}
	res := LockResult[*T](shepard.Ok[Guard[*T], error](guard))
	return &res
}

func (m *RWMutex[T]) RUnlock() {
	m.mutex.Unlock()
}

func (m *RWMutex[T]) Lock() *LockResult[*T] {
	m.mutex.Lock()
	guard := Guard[*T]{t: &m.t}
	res := LockResult[*T](shepard.Ok[Guard[*T], error](guard))
	return &res
}

func (m *RWMutex[T]) Unlock() {
	m.mutex.Unlock()
}
