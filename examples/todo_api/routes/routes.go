// Package routes configures all API routes.
package routes

import (
	"github.com/marlaone/shepard/examples/todo_api/handlers"
	"github.com/marlaone/shepard/examples/todo_api/middleware"
	"github.com/marlaone/shepard/http"
)

// Setup configures all routes on the router.
func Setup(router *http.Router, todoHandlers *handlers.TodoHandlers) {
	// API v1 group with logging middleware
	api := router.Group("/api/v1")
	api.Use(middleware.RequestLogger)

	// Todo CRUD routes
	api.Get("/todos", todoHandlers.List())
	api.Get("/todos/filter", todoHandlers.ListFiltered())
	api.Get("/todos/:id", todoHandlers.Get())
	api.Post("/todos", todoHandlers.Create())
	api.Patch("/todos/:id", todoHandlers.Update())
	api.Delete("/todos/:id", todoHandlers.Delete())
	api.Post("/todos/:id/complete", todoHandlers.Complete())

	// Stats route
	api.Get("/stats", todoHandlers.Stats())
}
