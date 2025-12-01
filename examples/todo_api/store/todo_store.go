// Package store provides data storage using Shepard's Rust-like collections.
package store

import (
	"time"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
	"github.com/marlaone/shepard/examples/todo_api/domain"
	"github.com/marlaone/shepard/sync/rwmutex"
)

// TodoStore provides thread-safe storage for todos using Shepard types:
// - rwmutex.RWMutex[T] instead of sync.RWMutex
// - hashmap.HashMap[K, V] instead of map[K]V
type TodoStore struct {
	todos  rwmutex.RWMutex[hashmap.HashMap[domain.TodoID, domain.Todo]]
	nextID rwmutex.RWMutex[domain.TodoID]
}

// New creates a new TodoStore with empty storage.
func New() *TodoStore {
	return &TodoStore{
		todos:  rwmutex.New(hashmap.New[domain.TodoID, domain.Todo]()),
		nextID: rwmutex.New(domain.TodoID(1)),
	}
}

// Create adds a new todo and returns it.
func (s *TodoStore) Create(req domain.CreateTodoRequest) domain.Todo {
	// Get and increment ID with write lock
	idGuard := s.nextID.Lock()
	idPtr := idGuard.Unwrap()
	currentID := *idPtr
	*idPtr = currentID + 1
	s.nextID.Unlock()

	now := time.Now()
	todo := domain.Todo{
		ID:          currentID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Insert into HashMap with write lock
	guard := s.todos.Lock()
	todos := guard.Unwrap()
	todos.Insert(todo.ID, todo)
	s.todos.Unlock()

	return todo
}

// Get retrieves a todo by ID, returns Option[Todo].
func (s *TodoStore) Get(id domain.TodoID) shepard.Option[domain.Todo] {
	guard := s.todos.RLock()
	todos := guard.Unwrap()
	defer s.todos.RUnlock()

	// HashMap.Get returns Option[*V]
	opt := todos.Get(id)
	if opt.IsNone() {
		return shepard.None[domain.Todo]()
	}
	return shepard.Some(*opt.Unwrap())
}

// List returns all todos as slice.Slice[Todo].
func (s *TodoStore) List() slice.Slice[domain.Todo] {
	guard := s.todos.RLock()
	todos := guard.Unwrap()
	defer s.todos.RUnlock()

	// Use HashMap's Values iterator with Shepard's Slice
	result := slice.WithCapacity[domain.Todo](todos.Len())
	todos.Values().Foreach(func(_ int, todo domain.Todo) {
		result.Push(todo)
	})
	return result
}

// Update modifies a todo, returns Option[Todo] (None if not found).
func (s *TodoStore) Update(id domain.TodoID, req domain.UpdateTodoRequest) shepard.Option[domain.Todo] {
	guard := s.todos.Lock()
	todos := guard.Unwrap()
	defer s.todos.Unlock()

	// Check if exists
	opt := todos.Get(id)
	if opt.IsNone() {
		return shepard.None[domain.Todo]()
	}

	existing := *opt.Unwrap()

	// Apply updates using Option pattern
	if req.Title.IsSome() {
		existing.Title = req.Title.Unwrap()
	}
	if req.Description.IsSome() {
		existing.Description = req.Description.Unwrap()
	}
	if req.Completed.IsSome() {
		existing.Completed = req.Completed.Unwrap()
	}
	existing.UpdatedAt = time.Now()

	// Update in HashMap
	todos.Insert(id, existing)
	return shepard.Some(existing)
}

// Delete removes a todo, returns Option[Todo] with removed item.
func (s *TodoStore) Delete(id domain.TodoID) shepard.Option[domain.Todo] {
	guard := s.todos.Lock()
	todos := guard.Unwrap()
	defer s.todos.Unlock()

	// HashMap.Remove returns Option[V]
	return todos.Remove(id)
}

// Count returns number of todos.
func (s *TodoStore) Count() int {
	guard := s.todos.RLock()
	todos := guard.Unwrap()
	defer s.todos.RUnlock()
	return todos.Len()
}
