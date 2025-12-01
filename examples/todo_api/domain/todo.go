// Package domain contains the core business types for the Todo API.
package domain

import (
	"time"

	"github.com/marlaone/shepard"
)

// TodoID is the unique identifier for a Todo item.
type TodoID int

// Todo represents a todo item in the system.
type Todo struct {
	ID          TodoID    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoRequest is the request body for creating a new todo.
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// UpdateTodoRequest is the request body for updating a todo.
// Uses Option types for partial updates.
type UpdateTodoRequest struct {
	Title       shepard.Option[string] `json:"title,omitempty"`
	Description shepard.Option[string] `json:"description,omitempty"`
	Completed   shepard.Option[bool]   `json:"completed,omitempty"`
}
