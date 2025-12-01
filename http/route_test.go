package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePattern_Static(t *testing.T) {
	result := ParsePattern("/users")

	assert.True(t, result.IsOk())
	pattern := result.Unwrap()

	assert.Equal(t, 1, pattern.Segments.Len())
	seg := pattern.Segments.Get(0)
	assert.True(t, seg.IsSome())
	assert.Equal(t, SegmentStatic, seg.Unwrap().Kind)
	assert.Equal(t, "users", seg.Unwrap().Value)
}

func TestParsePattern_WithParam(t *testing.T) {
	result := ParsePattern("/users/:id")

	assert.True(t, result.IsOk())
	pattern := result.Unwrap()

	assert.Equal(t, 2, pattern.Segments.Len())

	seg0 := pattern.Segments.Get(0)
	assert.True(t, seg0.IsSome())
	assert.Equal(t, SegmentStatic, seg0.Unwrap().Kind)
	assert.Equal(t, "users", seg0.Unwrap().Value)

	seg1 := pattern.Segments.Get(1)
	assert.True(t, seg1.IsSome())
	assert.Equal(t, SegmentParam, seg1.Unwrap().Kind)
	assert.Equal(t, "id", seg1.Unwrap().Value)
}

func TestParsePattern_MultipleParams(t *testing.T) {
	result := ParsePattern("/users/:userId/posts/:postId")

	assert.True(t, result.IsOk())
	pattern := result.Unwrap()

	assert.Equal(t, 4, pattern.Segments.Len())
}

func TestParsePattern_Wildcard(t *testing.T) {
	result := ParsePattern("/files/*path")

	assert.True(t, result.IsOk())
	pattern := result.Unwrap()

	seg1 := pattern.Segments.Get(1)
	assert.True(t, seg1.IsSome())
	assert.Equal(t, SegmentWildcard, seg1.Unwrap().Kind)
	assert.Equal(t, "path", seg1.Unwrap().Value)
}

func TestParsePattern_Root(t *testing.T) {
	result := ParsePattern("/")

	assert.True(t, result.IsOk())
	pattern := result.Unwrap()

	assert.Equal(t, 0, pattern.Segments.Len())
}

func TestRoutePattern_Match_Static(t *testing.T) {
	pattern := ParsePattern("/users").Unwrap()

	match := pattern.Match("/users")
	assert.True(t, match.IsSome())
	assert.Equal(t, 0, match.Unwrap().Len())

	noMatch := pattern.Match("/posts")
	assert.True(t, noMatch.IsNone())
}

func TestRoutePattern_Match_WithParam(t *testing.T) {
	pattern := ParsePattern("/users/:id").Unwrap()

	match := pattern.Match("/users/123")
	assert.True(t, match.IsSome())

	params := match.Unwrap()
	id := params.Get("id")
	assert.True(t, id.IsSome())
	assert.Equal(t, "123", *id.Unwrap())
}

func TestRoutePattern_Match_MultipleParams(t *testing.T) {
	pattern := ParsePattern("/users/:userId/posts/:postId").Unwrap()

	match := pattern.Match("/users/42/posts/100")
	assert.True(t, match.IsSome())

	params := match.Unwrap()

	userId := params.Get("userId")
	assert.True(t, userId.IsSome())
	assert.Equal(t, "42", *userId.Unwrap())

	postId := params.Get("postId")
	assert.True(t, postId.IsSome())
	assert.Equal(t, "100", *postId.Unwrap())
}

func TestRoutePattern_Match_WrongSegmentCount(t *testing.T) {
	pattern := ParsePattern("/users/:id").Unwrap()

	// Too few segments
	noMatch := pattern.Match("/users")
	assert.True(t, noMatch.IsNone())

	// Too many segments
	noMatch = pattern.Match("/users/123/extra")
	assert.True(t, noMatch.IsNone())
}

func TestRoutePattern_Match_Root(t *testing.T) {
	pattern := ParsePattern("/").Unwrap()

	match := pattern.Match("/")
	assert.True(t, match.IsSome())

	noMatch := pattern.Match("/users")
	assert.True(t, noMatch.IsNone())
}

func TestRoutePattern_Match_Wildcard(t *testing.T) {
	pattern := ParsePattern("/files/*path").Unwrap()

	match := pattern.Match("/files/a/b/c.txt")
	assert.True(t, match.IsSome())

	params := match.Unwrap()
	path := params.Get("path")
	assert.True(t, path.IsSome())
	assert.Equal(t, "a/b/c.txt", *path.Unwrap())
}

func TestRoutePattern_Original(t *testing.T) {
	pattern := ParsePattern("/users/:id").Unwrap()

	assert.Equal(t, "/users/:id", pattern.Original)
}

func TestRoute_New(t *testing.T) {
	handler := SimpleHandlerFunc(func(ctx *Context) Response {
		return Text("test")
	})

	// When there's no middleware, handler and chainedHandler are the same
	route := NewRoute(MethodGet, "/users/:id", handler, handler)

	assert.Equal(t, MethodGet, route.Method)
	assert.Equal(t, "/users/:id", route.Pattern.Original)
	assert.NotNil(t, route.Handler)
	assert.NotNil(t, route.ChainedHandler)
}
