package mutex

import (
	"fmt"
	"github.com/marlaone/shepard"
	"runtime"
	"sync"
)

type LockResult[T any] shepard.Result[Guard[T], error]

func (l *LockResult[T]) Unwrap() T {
	fmt.Println("unwrap")
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
	runtime.SetFinalizer(&m, func(g *Mutex[T]) {
		fmt.Println("mutex")
	})
	return m
}

func (m *Mutex[T]) Lock() *LockResult[T] {
	m.mutex.Lock()
	guard := Guard[T]{t: m.t}
	fmt.Println("lock")
	runtime.SetFinalizer(&guard, func(g *Guard[T]) {
		fmt.Println("unlock")
	})
	res := LockResult[T](shepard.Ok[Guard[T], error](guard))
	runtime.SetFinalizer(&res, func(res *LockResult[T]) {
		fmt.Println("res")
		m.unlock()
	})
	return &res
}

func (m *Mutex[T]) unlock() {
	m.mutex.Unlock()
}
