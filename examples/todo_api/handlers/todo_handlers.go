// Package handlers provides HTTP handlers for the Todo API.
package handlers

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/examples/todo_api/domain"
	"github.com/marlaone/shepard/examples/todo_api/store"
	"github.com/marlaone/shepard/http"
	"github.com/marlaone/shepard/http/extract"
)

// TodoHandlers contains all HTTP handlers for todo operations.
type TodoHandlers struct {
	store *store.TodoStore
}

// New creates a new TodoHandlers instance.
func New(store *store.TodoStore) *TodoHandlers {
	return &TodoHandlers{store: store}
}

// List handles GET /todos - returns all todos.
func (h *TodoHandlers) List() http.Handler {
	return http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		todos := h.store.List()
		// Convert slice.Slice to []T for JSON serialization
		return http.JsonBody(todos.Iter().Values()).IntoResponse()
	})
}

// Get handles GET /todos/:id - returns a single todo.
func (h *TodoHandlers) Get() http.Handler {
	return http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		// Extract path parameter with type safety
		id := extract.ExtractPath[int](ctx, "id")
		if id.IsErr() {
			return http.BadRequestErr("invalid todo ID")
		}

		// Get returns Option[Todo]
		opt := h.store.Get(domain.TodoID(id.Unwrap()))
		if opt.IsNone() {
			return http.NotFoundErr("todo not found")
		}

		return http.JsonOk(opt.Unwrap())
	})
}

// Create handles POST /todos - creates a new todo.
func (h *TodoHandlers) Create() http.Handler {
	return http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		// Extract JSON body
		body := extract.Json[domain.CreateTodoRequest](ctx)
		if body.IsErr() {
			return http.BadRequestErr("invalid request body")
		}

		req := body.Unwrap()

		// Validate
		if req.Title == "" {
			return http.BadRequestErr("title is required")
		}

		todo := h.store.Create(req)
		return http.JsonCreated(todo)
	})
}

// Update handles PATCH /todos/:id - updates a todo.
func (h *TodoHandlers) Update() http.Handler {
	return http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		id := extract.ExtractPath[int](ctx, "id")
		if id.IsErr() {
			return http.BadRequestErr("invalid todo ID")
		}

		body := extract.Json[domain.UpdateTodoRequest](ctx)
		if body.IsErr() {
			return http.BadRequestErr("invalid request body")
		}

		// Update returns Option[Todo]
		opt := h.store.Update(domain.TodoID(id.Unwrap()), body.Unwrap())
		if opt.IsNone() {
			return http.NotFoundErr("todo not found")
		}

		return http.JsonOk(opt.Unwrap())
	})
}

// Delete handles DELETE /todos/:id - deletes a todo.
func (h *TodoHandlers) Delete() http.Handler {
	return http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		id := extract.ExtractPath[int](ctx, "id")
		if id.IsErr() {
			return http.BadRequestErr("invalid todo ID")
		}

		// Delete returns Option[Todo]
		opt := h.store.Delete(domain.TodoID(id.Unwrap()))
		if opt.IsNone() {
			return http.NotFoundErr("todo not found")
		}

		return http.JsonOk(map[string]any{
			"message": "todo deleted",
			"deleted": opt.Unwrap(),
		})
	})
}

// Complete handles POST /todos/:id/complete - marks a todo as complete.
func (h *TodoHandlers) Complete() http.Handler {
	return http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		idResult := extract.ExtractPath[int](ctx, "id")
		if idResult.IsErr() {
			return http.BadRequestErr("invalid todo ID")
		}

		// Update with completed=true using Option
		updateReq := domain.UpdateTodoRequest{
			Completed: shepard.Some(true),
		}

		opt := h.store.Update(domain.TodoID(idResult.Unwrap()), updateReq)
		if opt.IsNone() {
			return http.NotFoundErr("todo not found")
		}

		return http.JsonOk(opt.Unwrap())
	})
}

// ListFiltered handles GET /todos/filter?completed=true&limit=10.
func (h *TodoHandlers) ListFiltered() http.Handler {
	return http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		// Optional query parameters with defaults
		completedFilter := extract.QueryOption[bool](ctx, "completed")
		limit := extract.QueryOr(ctx, "limit", 100)

		todos := h.store.List()

		// Filter using Option pattern with slice.Slice.FilterInPlace
		if completedFilter.IsSome() {
			filterVal := completedFilter.Unwrap()
			todos.Retain(func(t *domain.Todo) bool {
				return t.Completed == filterVal
			})
		}

		// Apply limit using slice.Slice.Truncate
		if todos.Len() > limit {
			todos.Truncate(uint(limit))
		}

		return http.JsonBody(todos).IntoResponse()
	})
}

// Stats handles GET /stats - returns statistics.
func (h *TodoHandlers) Stats() http.Handler {
	return http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.JsonBody(map[string]any{
			"total_todos": h.store.Count(),
		}).IntoResponse()
	})
}
