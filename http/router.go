package http

import (
	"github.com/marlaone/shepard"
)

// Router is the main HTTP router that matches requests to handlers
type Router struct {
	tree                    *RadixTree // Radix tree for O(k) route lookup
	middleware              []Middleware
	notFoundHandler         Handler
	chainedNotFound         Handler // Pre-computed 404 with middleware
	chainedMethodNotAllowed Handler // Pre-computed 405 with middleware
}

// NewRouter creates a new Router
func NewRouter() *Router {
	notFound := defaultNotFoundHandler()
	methodNotAllowed := defaultMethodNotAllowedHandler()
	return &Router{
		tree:                    NewRadixTree(),
		middleware:              nil,
		notFoundHandler:         notFound,
		chainedNotFound:         notFound,
		chainedMethodNotAllowed: methodNotAllowed,
	}
}

// buildChain wraps handler with all middleware (called at registration, not per-request)
func (r *Router) buildChain(handler Handler) Handler {
	result := handler
	// Apply in reverse order (last middleware = outermost)
	for i := len(r.middleware) - 1; i >= 0; i-- {
		result = r.middleware[i](result)
	}
	return result
}

// rebuildDefaultChains rebuilds pre-computed chains for 404/405 handlers
func (r *Router) rebuildDefaultChains() {
	r.chainedNotFound = r.buildChain(r.notFoundHandler)
	r.chainedMethodNotAllowed = r.buildChain(defaultMethodNotAllowedHandler())
}

// defaultNotFoundHandler returns the default 404 handler
func defaultNotFoundHandler() Handler {
	return SimpleHandlerFunc(func(ctx *Context) Response {
		return TextWithStatus(StatusNotFound, "404 Not Found")
	})
}

// defaultMethodNotAllowedHandler returns the default 405 handler
func defaultMethodNotAllowedHandler() Handler {
	return SimpleHandlerFunc(func(ctx *Context) Response {
		return TextWithStatus(StatusMethodNotAllowed, "405 Method Not Allowed")
	})
}

// Route adds a route with the given method, pattern, and handler
func (r *Router) Route(method Method, pattern string, handler Handler) *Router {
	// Pre-compute the middleware chain at registration time
	chainedHandler := r.buildChain(handler)

	// Insert into radix tree for O(k) lookup
	r.tree.Insert(method, pattern, handler, chainedHandler)

	return r
}

// Get adds a GET route
func (r *Router) Get(pattern string, handler Handler) *Router {
	return r.Route(MethodGet, pattern, handler)
}

// Post adds a POST route
func (r *Router) Post(pattern string, handler Handler) *Router {
	return r.Route(MethodPost, pattern, handler)
}

// Put adds a PUT route
func (r *Router) Put(pattern string, handler Handler) *Router {
	return r.Route(MethodPut, pattern, handler)
}

// Delete adds a DELETE route
func (r *Router) Delete(pattern string, handler Handler) *Router {
	return r.Route(MethodDelete, pattern, handler)
}

// Patch adds a PATCH route
func (r *Router) Patch(pattern string, handler Handler) *Router {
	return r.Route(MethodPatch, pattern, handler)
}

// Head adds a HEAD route
func (r *Router) Head(pattern string, handler Handler) *Router {
	return r.Route(MethodHead, pattern, handler)
}

// Options adds an OPTIONS route
func (r *Router) Options(pattern string, handler Handler) *Router {
	return r.Route(MethodOptions, pattern, handler)
}

// NotFound sets a custom 404 handler
func (r *Router) NotFound(handler Handler) *Router {
	r.notFoundHandler = handler
	r.rebuildDefaultChains()
	return r
}

// Use adds a middleware to the router
// NOTE: Middleware must be added BEFORE routes for them to be applied
func (r *Router) Use(mw Middleware) *Router {
	r.middleware = append(r.middleware, mw)
	r.rebuildDefaultChains()
	return r
}

// Handle processes a request through the router
func (r *Router) Handle(ctx *Context) shepard.Result[Response, HttpError] {
	// Find the pre-computed handler chain (no per-request allocations)
	handler := r.findHandler(ctx)

	// Execute the handler chain
	result := handler.Handle(ctx)

	// Convert HttpError to Response if needed
	if result.IsErr() {
		err := result.UnwrapErr()
		return shepard.Ok[Response, HttpError](err.IntoResponse())
	}

	return result
}

// findHandler finds the matching handler for a request
// Uses radix tree for O(k) lookup where k = path length
func (r *Router) findHandler(ctx *Context) Handler {
	path := ctx.Path()
	method := ctx.Method()

	// Radix tree lookup
	node, params := r.tree.Lookup(method, path)

	if node != nil && node.chainedHandler != nil {
		// Copy params to context
		for _, p := range params {
			ctx.SetPathParam(p.Name, p.Value)
		}
		return node.chainedHandler
	}

	// Check for 405 (path exists but wrong method)
	if r.tree.PathExistsAnyMethod(path) {
		return r.chainedMethodNotAllowed
	}

	return r.chainedNotFound
}

// Middleware wraps a Handler and returns a new Handler
// Pre-computes the chain at route registration time (zero per-request allocations)
type Middleware func(next Handler) Handler

// --- Route Groups ---

// RouteGroup represents a group of routes with a shared prefix and optional middleware
type RouteGroup struct {
	router     *Router
	prefix     string
	middleware []Middleware
}

// Group creates a new route group with the given prefix
// Routes added to the group will have the prefix prepended to their patterns
func (r *Router) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		router:     r,
		prefix:     prefix,
		middleware: nil,
	}
}

// Use adds middleware to this group only
// Middleware is applied to all routes added to this group
func (g *RouteGroup) Use(mw Middleware) *RouteGroup {
	g.middleware = append(g.middleware, mw)
	return g
}

// Route adds a route with the given method, pattern, and handler to the group
func (g *RouteGroup) Route(method Method, pattern string, handler Handler) *RouteGroup {
	g.addRoute(method, pattern, handler)
	return g
}

// Get adds a GET route to the group
func (g *RouteGroup) Get(pattern string, handler Handler) *RouteGroup {
	return g.Route(MethodGet, pattern, handler)
}

// Post adds a POST route to the group
func (g *RouteGroup) Post(pattern string, handler Handler) *RouteGroup {
	return g.Route(MethodPost, pattern, handler)
}

// Put adds a PUT route to the group
func (g *RouteGroup) Put(pattern string, handler Handler) *RouteGroup {
	return g.Route(MethodPut, pattern, handler)
}

// Delete adds a DELETE route to the group
func (g *RouteGroup) Delete(pattern string, handler Handler) *RouteGroup {
	return g.Route(MethodDelete, pattern, handler)
}

// Patch adds a PATCH route to the group
func (g *RouteGroup) Patch(pattern string, handler Handler) *RouteGroup {
	return g.Route(MethodPatch, pattern, handler)
}

// Head adds a HEAD route to the group
func (g *RouteGroup) Head(pattern string, handler Handler) *RouteGroup {
	return g.Route(MethodHead, pattern, handler)
}

// Options adds an OPTIONS route to the group
func (g *RouteGroup) Options(pattern string, handler Handler) *RouteGroup {
	return g.Route(MethodOptions, pattern, handler)
}

// Group creates a nested group with an additional prefix
func (g *RouteGroup) Group(prefix string) *RouteGroup {
	// Clone middleware slice
	mwCopy := make([]Middleware, len(g.middleware))
	copy(mwCopy, g.middleware)
	return &RouteGroup{
		router:     g.router,
		prefix:     g.prefix + prefix,
		middleware: mwCopy,
	}
}

// addRoute adds a route to the parent router with the group's prefix and middleware
func (g *RouteGroup) addRoute(method Method, pattern string, handler Handler) {
	fullPattern := g.prefix + pattern

	// Apply group middleware to handler (innermost to outermost)
	wrappedHandler := handler
	for i := len(g.middleware) - 1; i >= 0; i-- {
		wrappedHandler = g.middleware[i](wrappedHandler)
	}

	// Then router's buildChain will apply global middleware
	g.router.Route(method, fullPattern, wrappedHandler)
}
