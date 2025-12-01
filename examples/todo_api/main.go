// Package main demonstrates a real-world REST API using Shepard's Rust-like patterns.
//
// This Todo API showcases:
//   - hashmap.HashMap instead of Go maps
//   - slice.Slice instead of Go slices
//   - sync/rwmutex.RWMutex instead of sync.RWMutex
//   - shepard.Option for optional returns
//   - shepard.Result for error handling
//   - HTTP extractors and response helpers
//
// Project structure:
//
//	todo_api/
//	├── main.go           # Application entry point
//	├── domain/           # Business domain types
//	│   └── todo.go
//	├── store/            # Data storage layer
//	│   └── todo_store.go
//	├── handlers/         # HTTP handlers
//	│   └── todo_handlers.go
//	├── middleware/       # HTTP middleware
//	│   └── logging.go
//	└── routes/           # Route configuration
//	    └── routes.go
package main

import (
	"fmt"
	"log"

	"github.com/marlaone/shepard/examples/todo_api/domain"
	"github.com/marlaone/shepard/examples/todo_api/handlers"
	"github.com/marlaone/shepard/examples/todo_api/routes"
	"github.com/marlaone/shepard/examples/todo_api/store"
	"github.com/marlaone/shepard/http"
)

func main() {
	// Initialize store
	todoStore := store.New()

	// Seed some data
	seedData(todoStore)

	// Initialize handlers
	todoHandlers := handlers.New(todoStore)

	// Setup router and routes
	router := http.NewRouter()
	routes.Setup(router, todoHandlers)

	// Print banner
	printBanner()

	// Start server
	server := http.NewServer(":8080", router)
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}

func seedData(s *store.TodoStore) {
	s.Create(domain.CreateTodoRequest{
		Title:       "Learn Shepard",
		Description: "Explore Rust-like patterns in Go",
	})
	s.Create(domain.CreateTodoRequest{
		Title:       "Build REST API",
		Description: "Use HashMap, Option, Result",
	})
	s.Create(domain.CreateTodoRequest{
		Title: "Write tests",
	})
}

func printBanner() {
	fmt.Println("Shepard Todo API")
	fmt.Println("================")
	fmt.Println("Demonstrating Rust-like patterns:")
	fmt.Println("  - hashmap.HashMap[K, V] for storage")
	fmt.Println("  - slice.Slice[T] for collections")
	fmt.Println("  - rwmutex.RWMutex[T] for thread-safe access")
	fmt.Println("  - shepard.Option[T] for nullable values")
	fmt.Println("  - shepard.Result[T, E] for error handling")
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Println("  GET    /api/v1/todos              - List all todos")
	fmt.Println("  GET    /api/v1/todos/filter       - Filter todos (?completed=true&limit=10)")
	fmt.Println("  GET    /api/v1/todos/:id          - Get todo by ID")
	fmt.Println("  POST   /api/v1/todos              - Create todo")
	fmt.Println("  PATCH  /api/v1/todos/:id          - Update todo")
	fmt.Println("  DELETE /api/v1/todos/:id          - Delete todo")
	fmt.Println("  POST   /api/v1/todos/:id/complete - Mark complete")
	fmt.Println("  GET    /api/v1/stats              - Get stats")
	fmt.Println()
	fmt.Println("Starting server on :8080...")
}
