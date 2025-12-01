package http

import (
	"strings"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
)

// SegmentKind represents the type of a route segment
type SegmentKind uint8

const (
	SegmentStatic   SegmentKind = iota // Static segment like "users"
	SegmentParam                       // Parameter segment like ":id"
	SegmentWildcard                    // Wildcard segment like "*path"
)

// Segment represents a single part of a route pattern
type Segment struct {
	Kind  SegmentKind
	Value string
}

// RoutePattern represents a compiled route pattern
type RoutePattern struct {
	Segments slice.Slice[Segment]
	Original string
}

// ParsePattern parses a route pattern string into a RoutePattern
// Supports:
// - Static segments: /users
// - Parameters: /users/:id
// - Wildcards: /files/*path
func ParsePattern(pattern string) shepard.Result[RoutePattern, error] {
	// Handle root path
	if pattern == "/" {
		return shepard.Ok[RoutePattern, error](RoutePattern{
			Segments: slice.New[Segment](),
			Original: pattern,
		})
	}

	// Remove leading/trailing slashes and split
	trimmed := strings.Trim(pattern, "/")
	parts := strings.Split(trimmed, "/")

	segments := slice.WithCapacity[Segment](len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		var seg Segment
		if strings.HasPrefix(part, ":") {
			// Parameter segment
			seg = Segment{
				Kind:  SegmentParam,
				Value: strings.TrimPrefix(part, ":"),
			}
		} else if strings.HasPrefix(part, "*") {
			// Wildcard segment (must be last)
			seg = Segment{
				Kind:  SegmentWildcard,
				Value: strings.TrimPrefix(part, "*"),
			}
		} else {
			// Static segment
			seg = Segment{
				Kind:  SegmentStatic,
				Value: part,
			}
		}
		segments.Push(seg)
	}

	return shepard.Ok[RoutePattern, error](RoutePattern{
		Segments: segments,
		Original: pattern,
	})
}

// ParamCount returns the number of parameter segments in the pattern
func (rp RoutePattern) ParamCount() int {
	count := 0
	for i := 0; i < rp.Segments.Len(); i++ {
		if seg := rp.Segments.Get(i); seg.IsSome() {
			s := seg.Unwrap()
			if s.Kind == SegmentParam || s.Kind == SegmentWildcard {
				count++
			}
		}
	}
	return count
}

// Match attempts to match a path against this pattern
// Returns Option[HashMap[string, string]] with extracted parameters if matched
// Uses index-based iteration to avoid allocations from strings.Split
func (rp RoutePattern) Match(path string) shepard.Option[hashmap.HashMap[string, string]] {
	// Handle root path
	if path == "/" || path == "" {
		if rp.Segments.Len() == 0 {
			return shepard.Some(hashmap.New[string, string]())
		}
		return shepard.None[hashmap.HashMap[string, string]]()
	}

	// Check for wildcard at the end
	hasWildcard := false
	wildcardName := ""
	numSegments := rp.Segments.Len()
	if numSegments > 0 {
		lastSeg := rp.Segments.Get(numSegments - 1)
		if lastSeg.IsSome() && lastSeg.Unwrap().Kind == SegmentWildcard {
			hasWildcard = true
			wildcardName = lastSeg.Unwrap().Value
		}
	}

	segmentsToMatch := numSegments
	if hasWildcard {
		segmentsToMatch--
	}

	// Trim leading slash for iteration
	pathStart := 0
	pathEnd := len(path)
	if pathEnd > 0 && path[0] == '/' {
		pathStart = 1
	}
	if pathEnd > pathStart && path[pathEnd-1] == '/' {
		pathEnd--
	}

	// Count path segments and validate length match (without allocating)
	pathSegmentCount := 0
	if pathEnd > pathStart {
		pathSegmentCount = 1
		for i := pathStart; i < pathEnd; i++ {
			if path[i] == '/' {
				pathSegmentCount++
			}
		}
	}

	// For non-wildcard patterns, segment count must match
	if !hasWildcard && numSegments != pathSegmentCount {
		return shepard.None[hashmap.HashMap[string, string]]()
	}

	// For wildcard patterns, we need at least (segments - 1) path parts
	if hasWildcard && pathSegmentCount < segmentsToMatch {
		return shepard.None[hashmap.HashMap[string, string]]()
	}

	// Temporary storage for param values (to defer HashMap allocation until match confirmed)
	// We use a fixed-size array to avoid allocation; most routes have few params
	type paramPair struct {
		name  string
		value string
	}
	var paramStorage [8]paramPair
	paramCount := 0

	// Index-based iteration through path segments
	segmentIdx := 0
	pos := pathStart

	for segmentIdx < segmentsToMatch && pos < pathEnd {
		// Find end of current path segment
		segEnd := pos
		for segEnd < pathEnd && path[segEnd] != '/' {
			segEnd++
		}
		pathPart := path[pos:segEnd]

		// Get pattern segment
		segOpt := rp.Segments.Get(segmentIdx)
		if segOpt.IsNone() {
			return shepard.None[hashmap.HashMap[string, string]]()
		}
		seg := segOpt.Unwrap()

		switch seg.Kind {
		case SegmentStatic:
			if seg.Value != pathPart {
				return shepard.None[hashmap.HashMap[string, string]]()
			}
		case SegmentParam:
			if paramCount < len(paramStorage) {
				paramStorage[paramCount] = paramPair{name: seg.Value, value: pathPart}
				paramCount++
			}
		}

		segmentIdx++
		pos = segEnd + 1 // Skip past the '/'
	}

	// Verify we matched all required segments
	if segmentIdx < segmentsToMatch {
		return shepard.None[hashmap.HashMap[string, string]]()
	}

	// For non-wildcard, verify no extra path segments remain
	if !hasWildcard && pos < pathEnd {
		return shepard.None[hashmap.HashMap[string, string]]()
	}

	// Match succeeded - now allocate params HashMap
	params := hashmap.WithCapacity[string, string](paramCount + 1)

	// Copy stored params to HashMap
	for i := 0; i < paramCount; i++ {
		params.Insert(paramStorage[i].name, paramStorage[i].value)
	}

	// Handle wildcard - capture remaining path
	if hasWildcard {
		remaining := ""
		if pos < pathEnd {
			remaining = path[pos:pathEnd]
		} else if pos == pathEnd && pathEnd > pathStart {
			// Edge case: wildcard matches empty string when path ends exactly at wildcard position
			remaining = ""
		}
		params.Insert(wildcardName, remaining)
	}

	return shepard.Some(params)
}

// Route represents a single route with method, pattern, and handler
type Route struct {
	Method         Method
	Pattern        RoutePattern
	Handler        Handler // Original handler
	ChainedHandler Handler // Pre-computed: middleware chain + handler
}

// NewRoute creates a new Route
// chainedHandler should be the handler with middleware pre-applied
func NewRoute(method Method, pattern string, handler Handler, chainedHandler Handler) Route {
	patternResult := ParsePattern(pattern)
	pat := RoutePattern{Original: pattern, Segments: slice.New[Segment]()}
	if patternResult.IsOk() {
		pat = patternResult.Unwrap()
	}

	return Route{
		Method:         method,
		Pattern:        pat,
		Handler:        handler,
		ChainedHandler: chainedHandler,
	}
}

// Match attempts to match a method and path against this route
// Returns Option[HashMap[string, string]] with extracted parameters if matched
func (r Route) Match(method Method, path string) shepard.Option[hashmap.HashMap[string, string]] {
	if r.Method != method {
		return shepard.None[hashmap.HashMap[string, string]]()
	}
	return r.Pattern.Match(path)
}
