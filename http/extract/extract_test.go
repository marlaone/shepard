package extract

import (
	"testing"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/http"
	"github.com/stretchr/testify/assert"
)

func TestPath_Extract(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/users/123", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetPathParam("id", "123")

	var path Path[string]
	result := path.Extract(ctx, "id")

	assert.True(t, result.IsOk())
	assert.Equal(t, "123", result.Unwrap())
}

func TestPath_ExtractInt(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/users/42", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetPathParam("id", "42")

	var path Path[int]
	result := path.Extract(ctx, "id")

	assert.True(t, result.IsOk())
	assert.Equal(t, 42, result.Unwrap())
}

func TestPath_NotFound(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/users", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	var path Path[string]
	result := path.Extract(ctx, "id")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestQuery_Extract(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/search?q=hello", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	var query Query[string]
	result := query.Extract(ctx, "q")

	assert.True(t, result.IsOk())
	assert.Equal(t, "hello", result.Unwrap())
}

func TestQuery_ExtractInt(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/page?num=5", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	var query Query[int]
	result := query.Extract(ctx, "num")

	assert.True(t, result.IsOk())
	assert.Equal(t, 5, result.Unwrap())
}

func TestQuery_NotFound(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/search", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	var query Query[string]
	result := query.Extract(ctx, "q")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestHeader_Extract(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	req.Headers.Set("Authorization", "Bearer token123")
	ctx := http.NewContext(req)

	var header Header
	result := header.Extract(ctx, "Authorization")

	assert.True(t, result.IsOk())
	assert.Equal(t, "Bearer token123", result.Unwrap())
}

func TestHeader_NotFound(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	var header Header
	result := header.Extract(ctx, "Authorization")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestJson_Extract(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	body := `{"name":"John","age":30}`
	req := http.NewRequest(http.MethodPost, "/users", http.NewBodyReaderFromBytes([]byte(body)))
	req.Headers.Set("Content-Type", "application/json")
	ctx := http.NewContext(req)

	result := Json[User](ctx)

	assert.True(t, result.IsOk())
	user := result.Unwrap()
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, 30, user.Age)
}

func TestJson_InvalidJson(t *testing.T) {
	body := `{invalid json}`
	req := http.NewRequest(http.MethodPost, "/users", http.NewBodyReaderFromBytes([]byte(body)))
	req.Headers.Set("Content-Type", "application/json")
	ctx := http.NewContext(req)

	type User struct {
		Name string `json:"name"`
	}

	result := Json[User](ctx)

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestState_Extract(t *testing.T) {
	type AppState struct {
		DBConn string
	}

	state := &AppState{DBConn: "postgres://localhost"}
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetState("app", state)

	var stateExtractor State[*AppState]
	result := stateExtractor.Extract(ctx, "app")

	assert.True(t, result.IsOk())
	assert.Equal(t, "postgres://localhost", result.Unwrap().DBConn)
}

func TestState_NotFound(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	var stateExtractor State[*struct{}]
	result := stateExtractor.Extract(ctx, "app")

	assert.True(t, result.IsErr())
}

// Test optional extractors that return Option instead of Result
func TestOptionalQuery(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/search?q=hello", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	result := OptionalQuery[string](ctx, "q")
	assert.True(t, result.IsSome())
	assert.Equal(t, "hello", result.Unwrap())

	// Missing query param returns None
	result = OptionalQuery[string](ctx, "missing")
	assert.True(t, result.IsNone())
}

func TestOptionalPath(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/users/123", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetPathParam("id", "123")

	result := OptionalPath[string](ctx, "id")
	assert.True(t, result.IsSome())
	assert.Equal(t, "123", result.Unwrap())

	// Missing path param returns None
	result = OptionalPath[string](ctx, "missing")
	assert.True(t, result.IsNone())
}

// Test Extractor helper function for generic extraction
func TestExtractPath(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/users/42", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetPathParam("id", "42")

	result := ExtractPath[int](ctx, "id")

	assert.True(t, result.IsOk())
	assert.Equal(t, 42, result.Unwrap())
}

func TestExtractQuery(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/search?page=5", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	result := ExtractQuery[int](ctx, "page")

	assert.True(t, result.IsOk())
	assert.Equal(t, 5, result.Unwrap())
}

// Test extracting multiple path params
type UserParams struct {
	OrgID  string
	UserID int
}

func TestExtractMultipleParams(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/orgs/acme/users/42", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetPathParam("org_id", "acme")
	ctx.SetPathParam("user_id", "42")

	// Extract org_id as string
	orgResult := ExtractPath[string](ctx, "org_id")
	assert.True(t, orgResult.IsOk())
	assert.Equal(t, "acme", orgResult.Unwrap())

	// Extract user_id as int
	userResult := ExtractPath[int](ctx, "user_id")
	assert.True(t, userResult.IsOk())
	assert.Equal(t, 42, userResult.Unwrap())
}

// Test parsing errors for numeric types
func TestPath_ParseError(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/users/abc", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetPathParam("id", "abc")

	var path Path[int]
	result := path.Extract(ctx, "id")

	assert.True(t, result.IsErr())
	err := result.UnwrapErr()
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.Contains(t, err.Message, "parse")
}

func TestMethodExtract(t *testing.T) {
	req := http.NewRequest(http.MethodPost, "/users", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	method := ExtractMethod(ctx)
	assert.Equal(t, http.MethodPost, method)
}

// Helper for extracting multiple values at once using Result combinators
func TestCombineResults(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/orgs/acme/users/42?active=true", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	ctx.SetPathParam("org_id", "acme")
	ctx.SetPathParam("user_id", "42")

	// Use AndThen to chain extractions
	orgResult := ExtractPath[string](ctx, "org_id")
	userResult := ExtractPath[int](ctx, "user_id")
	activeResult := ExtractQuery[bool](ctx, "active")

	// All should succeed
	assert.True(t, orgResult.IsOk())
	assert.True(t, userResult.IsOk())
	assert.True(t, activeResult.IsOk())

	assert.Equal(t, "acme", orgResult.Unwrap())
	assert.Equal(t, 42, userResult.Unwrap())
	assert.Equal(t, true, activeResult.Unwrap())
}

// Test for extracting bool query params
func TestQuery_ExtractBool(t *testing.T) {
	tests := []struct {
		query    string
		expected bool
	}{
		{"active=true", true},
		{"active=false", false},
		{"active=1", true},
		{"active=0", false},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			req := http.NewRequest(http.MethodGet, "/users?"+tt.query, http.EmptyBodyReader())
			ctx := http.NewContext(req)

			result := ExtractQuery[bool](ctx, "active")
			assert.True(t, result.IsOk())
			assert.Equal(t, tt.expected, result.Unwrap())
		})
	}
}

// Test extracting float values
func TestQuery_ExtractFloat(t *testing.T) {
	req := http.NewRequest(http.MethodGet, "/price?amount=19.99", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	result := ExtractQuery[float64](ctx, "amount")
	assert.True(t, result.IsOk())
	assert.InDelta(t, 19.99, result.Unwrap(), 0.001)
}

// FromRequest interface test
func TestFromRequest(t *testing.T) {
	type CustomData struct {
		Value string
	}

	// Register a custom extractor
	RegisterExtractor("custom", func(ctx *http.Context) shepard.Result[any, http.HttpError] {
		return shepard.Ok[any, http.HttpError](CustomData{Value: "extracted"})
	})

	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	result := ExtractCustom[CustomData](ctx, "custom")
	assert.True(t, result.IsOk())
	assert.Equal(t, "extracted", result.Unwrap().Value)
}
