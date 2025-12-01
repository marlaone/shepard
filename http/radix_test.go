package http

import (
	"testing"
)

func TestRadixTree_Insert_Static(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/users", handler, handler)

	node, params := tree.Lookup(MethodGet, "/users")
	if node == nil {
		t.Fatal("expected to find node")
	}
	if node.chainedHandler == nil {
		t.Fatal("expected handler")
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestRadixTree_Insert_RootPath(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("root") })

	tree.Insert(MethodGet, "/", handler, handler)

	node, params := tree.Lookup(MethodGet, "/")
	if node == nil {
		t.Fatal("expected to find node for root path")
	}
	if node.chainedHandler == nil {
		t.Fatal("expected handler")
	}
	if len(params) != 0 {
		t.Errorf("expected 0 params, got %d", len(params))
	}
}

func TestRadixTree_Insert_DeepPath(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/api/v1/users/settings", handler, handler)

	node, _ := tree.Lookup(MethodGet, "/api/v1/users/settings")
	if node == nil {
		t.Fatal("expected to find node")
	}
}

func TestRadixTree_Lookup_Param(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/users/:id", handler, handler)

	node, params := tree.Lookup(MethodGet, "/users/123")
	if node == nil {
		t.Fatal("expected to find node")
	}
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	if params[0].Name != "id" {
		t.Errorf("expected param name 'id', got '%s'", params[0].Name)
	}
	if params[0].Value != "123" {
		t.Errorf("expected param value '123', got '%s'", params[0].Value)
	}
}

func TestRadixTree_Lookup_MultipleParams(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/users/:userId/posts/:postId", handler, handler)

	node, params := tree.Lookup(MethodGet, "/users/42/posts/99")
	if node == nil {
		t.Fatal("expected to find node")
	}
	if len(params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(params))
	}
	if params[0].Name != "userId" || params[0].Value != "42" {
		t.Errorf("first param: expected userId=42, got %s=%s", params[0].Name, params[0].Value)
	}
	if params[1].Name != "postId" || params[1].Value != "99" {
		t.Errorf("second param: expected postId=99, got %s=%s", params[1].Name, params[1].Value)
	}
}

func TestRadixTree_Lookup_Wildcard(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/files/*path", handler, handler)

	node, params := tree.Lookup(MethodGet, "/files/docs/readme.md")
	if node == nil {
		t.Fatal("expected to find node")
	}
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	if params[0].Name != "path" {
		t.Errorf("expected param name 'path', got '%s'", params[0].Name)
	}
	if params[0].Value != "docs/readme.md" {
		t.Errorf("expected 'docs/readme.md', got '%s'", params[0].Value)
	}
}

func TestRadixTree_StaticPriorityOverParam(t *testing.T) {
	tree := NewRadixTree()
	staticHandler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("static") })
	paramHandler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("param") })

	// Order of insertion shouldn't matter
	tree.Insert(MethodGet, "/users/:id", paramHandler, paramHandler)
	tree.Insert(MethodGet, "/users/new", staticHandler, staticHandler)

	// Static should win
	node, params := tree.Lookup(MethodGet, "/users/new")
	if node == nil {
		t.Fatal("expected to find node")
	}
	if len(params) != 0 {
		t.Error("static route should have no params")
	}

	// Param should still work
	node, params = tree.Lookup(MethodGet, "/users/123")
	if node == nil {
		t.Fatal("expected to find node")
	}
	if len(params) != 1 {
		t.Error("param route should have 1 param")
	}
}

func TestRadixTree_MethodSeparation(t *testing.T) {
	tree := NewRadixTree()
	getHandler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("get") })
	postHandler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("post") })

	tree.Insert(MethodGet, "/users", getHandler, getHandler)
	tree.Insert(MethodPost, "/users", postHandler, postHandler)

	// GET should work
	node, _ := tree.Lookup(MethodGet, "/users")
	if node == nil {
		t.Fatal("expected GET handler")
	}

	// POST should work
	node, _ = tree.Lookup(MethodPost, "/users")
	if node == nil {
		t.Fatal("expected POST handler")
	}

	// DELETE should not work
	node, _ = tree.Lookup(MethodDelete, "/users")
	if node != nil {
		t.Error("expected no DELETE handler")
	}
}

func TestRadixTree_NotFound(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/users", handler, handler)

	// Non-existent path
	node, _ := tree.Lookup(MethodGet, "/posts")
	if node != nil {
		t.Error("expected nil for non-existent path")
	}

	// Wrong method
	node, _ = tree.Lookup(MethodPost, "/users")
	if node != nil {
		t.Error("expected nil for wrong method")
	}
}

func TestRadixTree_PathExistsAnyMethod(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/users", handler, handler)
	tree.Insert(MethodPost, "/posts", handler, handler)

	if !tree.PathExistsAnyMethod("/users") {
		t.Error("expected /users to exist")
	}
	if !tree.PathExistsAnyMethod("/posts") {
		t.Error("expected /posts to exist")
	}
	if tree.PathExistsAnyMethod("/comments") {
		t.Error("expected /comments to not exist")
	}
}

func TestRadixTree_TrailingSlash(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/users", handler, handler)

	// With trailing slash should still match
	node, _ := tree.Lookup(MethodGet, "/users/")
	if node == nil {
		t.Fatal("expected to find node with trailing slash")
	}
}

func TestRadixTree_EmptyPath(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/", handler, handler)

	// Empty string should match root
	node, _ := tree.Lookup(MethodGet, "")
	if node == nil {
		t.Fatal("expected to find root node with empty path")
	}
}

func TestRadixTree_WildcardEmpty(t *testing.T) {
	tree := NewRadixTree()
	handler := SimpleHandlerFunc(func(ctx *Context) Response { return Text("ok") })

	tree.Insert(MethodGet, "/files/*path", handler, handler)

	// Wildcard with just the prefix
	node, params := tree.Lookup(MethodGet, "/files/")
	if node == nil {
		t.Fatal("expected to find node")
	}
	if len(params) != 1 {
		t.Fatalf("expected 1 param, got %d", len(params))
	}
	// Empty wildcard value
	if params[0].Value != "" {
		t.Errorf("expected empty wildcard value, got '%s'", params[0].Value)
	}
}

func TestRouter_RadixTree_Integration(t *testing.T) {
	router := NewRouter()

	router.Get("/", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("root")
	}))

	router.Get("/users/:id", SimpleHandlerFunc(func(ctx *Context) Response {
		id := ctx.Param("id")
		return Text("user: " + id.UnwrapOr(""))
	}))

	// Test root
	req := NewRequest(MethodGet, "/", EmptyBodyReader())
	ctx := NewContext(req)
	result := router.Handle(ctx)
	if result.IsErr() {
		t.Fatal("expected OK result")
	}

	// Test param
	req = NewRequest(MethodGet, "/users/42", EmptyBodyReader())
	ctx = NewContext(req)
	result = router.Handle(ctx)
	if result.IsErr() {
		t.Fatal("expected OK result")
	}
	resp := result.Unwrap()
	if resp.StatusCode != StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRouter_RadixTree_404(t *testing.T) {
	router := NewRouter()
	router.Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("ok")
	}))

	req := NewRequest(MethodGet, "/nonexistent", EmptyBodyReader())
	ctx := NewContext(req)
	result := router.Handle(ctx)
	if result.IsErr() {
		t.Fatal("expected OK result")
	}
	resp := result.Unwrap()
	if resp.StatusCode != StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestRouter_RadixTree_405(t *testing.T) {
	router := NewRouter()
	router.Get("/users", SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("ok")
	}))

	// POST to GET-only route
	req := NewRequest(MethodPost, "/users", EmptyBodyReader())
	ctx := NewContext(req)
	result := router.Handle(ctx)
	if result.IsErr() {
		t.Fatal("expected OK result")
	}
	resp := result.Unwrap()
	if resp.StatusCode != StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}
