package extract

import (
	"strings"
	"testing"

	"github.com/marlaone/shepard/http"
)

// Helper to create a context with path params
func createContextWithPathParam(key, value string) *http.Context {
	req := http.NewRequest(http.MethodGet, "/test", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetPathParam(key, value)
	return ctx
}

// Helper to create a context with query params
func createContextWithQuery(path string) *http.Context {
	req := http.NewRequest(http.MethodGet, path, http.EmptyBodyReader())
	return http.NewContext(req)
}

// Helper to create a context with headers
func createContextWithHeader(key, value string) *http.Context {
	req := http.NewRequest(http.MethodGet, "/test", http.EmptyBodyReader())
	req.Headers.Set(key, value)
	return http.NewContext(req)
}

// Helper to create a context with JSON body
func createContextWithJsonBody(jsonStr string) *http.Context {
	body := http.NewBodyReaderFromBytes([]byte(jsonStr))
	req := http.NewRequest(http.MethodPost, "/test", body)
	req.Headers.Set("Content-Type", "application/json")
	return http.NewContext(req)
}

// --- MustPath Tests ---

func TestMustPath_Success(t *testing.T) {
	ctx := createContextWithPathParam("id", "42")

	id := MustPath[int](ctx, "id")
	if id != 42 {
		t.Errorf("expected 42, got %d", id)
	}
}

func TestMustPath_String(t *testing.T) {
	ctx := createContextWithPathParam("name", "john")

	name := MustPath[string](ctx, "name")
	if name != "john" {
		t.Errorf("expected john, got %s", name)
	}
}

func TestMustPath_Panics_Missing(t *testing.T) {
	ctx := createContextWithPathParam("other", "value")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing path param")
		}
	}()

	MustPath[int](ctx, "id")
}

func TestMustPath_Panics_ParseError(t *testing.T) {
	ctx := createContextWithPathParam("id", "notanumber")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for parse error")
		}
	}()

	MustPath[int](ctx, "id")
}

// --- MustQuery Tests ---

func TestMustQuery_Success(t *testing.T) {
	ctx := createContextWithQuery("/test?page=5")

	page := MustQuery[int](ctx, "page")
	if page != 5 {
		t.Errorf("expected 5, got %d", page)
	}
}

func TestMustQuery_Panics_Missing(t *testing.T) {
	ctx := createContextWithQuery("/test")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing query param")
		}
	}()

	MustQuery[int](ctx, "page")
}

// --- MustHeader Tests ---

func TestMustHeader_Success(t *testing.T) {
	ctx := createContextWithHeader("Authorization", "Bearer token123")

	auth := MustHeader(ctx, "Authorization")
	if auth != "Bearer token123" {
		t.Errorf("expected Bearer token123, got %s", auth)
	}
}

func TestMustHeader_Panics_Missing(t *testing.T) {
	ctx := createContextWithHeader("Other", "value")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing header")
		}
	}()

	MustHeader(ctx, "Authorization")
}

// --- MustJson Tests ---

type TestUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestMustJson_Success(t *testing.T) {
	ctx := createContextWithJsonBody(`{"name": "Alice", "age": 30}`)

	user := MustJson[TestUser](ctx)
	if user.Name != "Alice" {
		t.Errorf("expected Alice, got %s", user.Name)
	}
	if user.Age != 30 {
		t.Errorf("expected 30, got %d", user.Age)
	}
}

func TestMustJson_Panics_InvalidJson(t *testing.T) {
	ctx := createContextWithJsonBody(`{invalid json}`)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid JSON")
		}
	}()

	MustJson[TestUser](ctx)
}

// --- PathOr Tests ---

func TestPathOr_Found(t *testing.T) {
	ctx := createContextWithPathParam("id", "42")

	id := PathOr(ctx, "id", 0)
	if id != 42 {
		t.Errorf("expected 42, got %d", id)
	}
}

func TestPathOr_NotFound_ReturnsDefault(t *testing.T) {
	ctx := createContextWithPathParam("other", "value")

	id := PathOr(ctx, "id", 100)
	if id != 100 {
		t.Errorf("expected default 100, got %d", id)
	}
}

func TestPathOr_ParseError_ReturnsDefault(t *testing.T) {
	ctx := createContextWithPathParam("id", "notanumber")

	id := PathOr(ctx, "id", 50)
	if id != 50 {
		t.Errorf("expected default 50, got %d", id)
	}
}

// --- QueryOr Tests ---

func TestQueryOr_Found(t *testing.T) {
	ctx := createContextWithQuery("/test?limit=25")

	limit := QueryOr(ctx, "limit", 10)
	if limit != 25 {
		t.Errorf("expected 25, got %d", limit)
	}
}

func TestQueryOr_NotFound_ReturnsDefault(t *testing.T) {
	ctx := createContextWithQuery("/test")

	limit := QueryOr(ctx, "limit", 10)
	if limit != 10 {
		t.Errorf("expected default 10, got %d", limit)
	}
}

// --- HeaderOr Tests ---

func TestHeaderOr_Found(t *testing.T) {
	ctx := createContextWithHeader("X-Request-ID", "abc123")

	reqID := HeaderOr(ctx, "X-Request-ID", "default")
	if reqID != "abc123" {
		t.Errorf("expected abc123, got %s", reqID)
	}
}

func TestHeaderOr_NotFound_ReturnsDefault(t *testing.T) {
	ctx := createContextWithHeader("Other", "value")

	reqID := HeaderOr(ctx, "X-Request-ID", "default-id")
	if reqID != "default-id" {
		t.Errorf("expected default-id, got %s", reqID)
	}
}

// --- TryPath Tests ---

func TestTryPath_Found(t *testing.T) {
	ctx := createContextWithPathParam("id", "42")

	id, ok := TryPath[int](ctx, "id")
	if !ok {
		t.Error("expected ok=true")
	}
	if id != 42 {
		t.Errorf("expected 42, got %d", id)
	}
}

func TestTryPath_NotFound(t *testing.T) {
	ctx := createContextWithPathParam("other", "value")

	id, ok := TryPath[int](ctx, "id")
	if ok {
		t.Error("expected ok=false")
	}
	if id != 0 {
		t.Errorf("expected zero value, got %d", id)
	}
}

// --- TryQuery Tests ---

func TestTryQuery_Found(t *testing.T) {
	ctx := createContextWithQuery("/test?active=true")

	active, ok := TryQuery[bool](ctx, "active")
	if !ok {
		t.Error("expected ok=true")
	}
	if !active {
		t.Error("expected true")
	}
}

func TestTryQuery_NotFound(t *testing.T) {
	ctx := createContextWithQuery("/test")

	_, ok := TryQuery[bool](ctx, "active")
	if ok {
		t.Error("expected ok=false")
	}
}

// --- TryHeader Tests ---

func TestTryHeader_Found(t *testing.T) {
	ctx := createContextWithHeader("Content-Type", "application/json")

	ct, ok := TryHeader(ctx, "Content-Type")
	if !ok {
		t.Error("expected ok=true")
	}
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestTryHeader_NotFound(t *testing.T) {
	ctx := createContextWithHeader("Other", "value")

	_, ok := TryHeader(ctx, "Content-Type")
	if ok {
		t.Error("expected ok=false")
	}
}

// --- TryJson Tests ---

func TestTryJson_Success(t *testing.T) {
	ctx := createContextWithJsonBody(`{"name": "Bob", "age": 25}`)

	user, ok := TryJson[TestUser](ctx)
	if !ok {
		t.Error("expected ok=true")
	}
	if user.Name != "Bob" {
		t.Errorf("expected Bob, got %s", user.Name)
	}
}

func TestTryJson_Failure(t *testing.T) {
	ctx := createContextWithJsonBody(`invalid`)

	_, ok := TryJson[TestUser](ctx)
	if ok {
		t.Error("expected ok=false for invalid JSON")
	}
}

// --- PathOption Tests ---

func TestPathOption_Found(t *testing.T) {
	ctx := createContextWithPathParam("id", "42")

	opt := PathOption[int](ctx, "id")
	if opt.IsNone() {
		t.Fatal("expected Some")
	}
	if opt.Unwrap() != 42 {
		t.Errorf("expected 42, got %d", opt.Unwrap())
	}
}

func TestPathOption_NotFound(t *testing.T) {
	ctx := createContextWithPathParam("other", "value")

	opt := PathOption[int](ctx, "id")
	if opt.IsSome() {
		t.Error("expected None")
	}
}

// --- QueryOption Tests ---

func TestQueryOption_Found(t *testing.T) {
	ctx := createContextWithQuery("/test?name=alice")

	opt := QueryOption[string](ctx, "name")
	if opt.IsNone() {
		t.Fatal("expected Some")
	}
	if opt.Unwrap() != "alice" {
		t.Errorf("expected alice, got %s", opt.Unwrap())
	}
}

func TestQueryOption_NotFound(t *testing.T) {
	ctx := createContextWithQuery("/test")

	opt := QueryOption[string](ctx, "name")
	if opt.IsSome() {
		t.Error("expected None")
	}
}

// --- HeaderOption Tests ---

func TestHeaderOption_Found(t *testing.T) {
	ctx := createContextWithHeader("Accept", "text/html")

	opt := HeaderOption(ctx, "Accept")
	if opt.IsNone() {
		t.Fatal("expected Some")
	}
	if opt.Unwrap() != "text/html" {
		t.Errorf("expected text/html, got %s", opt.Unwrap())
	}
}

func TestHeaderOption_NotFound(t *testing.T) {
	ctx := createContextWithHeader("Other", "value")

	opt := HeaderOption(ctx, "Accept")
	if opt.IsSome() {
		t.Error("expected None")
	}
}

// --- Integration Tests ---

func TestExtractors_RealWorldHandler(t *testing.T) {
	// Simulate a real handler using multiple extractors
	body := http.NewBodyReaderFromBytes([]byte(`{"name": "Test User"}`))
	req := http.NewRequest(http.MethodPost, "/users?notify=true", body)
	req.Headers.Set("Authorization", "Bearer secret")
	req.Headers.Set("Content-Type", "application/json")

	ctx := http.NewContext(req)
	ctx.SetPathParam("org", "acme")

	// Extract all values
	org := MustPath[string](ctx, "org")
	if org != "acme" {
		t.Errorf("expected acme, got %s", org)
	}

	notify := QueryOr(ctx, "notify", false)
	if !notify {
		t.Error("expected notify=true")
	}

	auth := MustHeader(ctx, "Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		t.Error("expected Bearer token")
	}

	type CreateUserRequest struct {
		Name string `json:"name"`
	}
	payload := MustJson[CreateUserRequest](ctx)
	if payload.Name != "Test User" {
		t.Errorf("expected Test User, got %s", payload.Name)
	}
}
