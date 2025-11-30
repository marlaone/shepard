package http

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/slice"
)

// Router is the main HTTP router that matches requests to handlers
type Router struct {
	routes          slice.Slice[Route]
	middleware      slice.Slice[Middleware]
	notFoundHandler Handler
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return &Router{
		routes:          slice.New[Route](),
		middleware:      slice.New[Middleware](),
		notFoundHandler: defaultNotFoundHandler(),
	}
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
	route := NewRoute(method, pattern, handler)
	r.routes.Push(route)
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
	return r
}

// Use adds a middleware to the router
func (r *Router) Use(mw Middleware) *Router {
	r.middleware.Push(mw)
	return r
}

// Handle processes a request through the router
func (r *Router) Handle(ctx *Context) shepard.Result[Response, HttpError] {
	// Build the handler chain with middleware
	handler := r.findHandler(ctx)

	// Apply middleware in reverse order (last added = outermost)
	for i := r.middleware.Len() - 1; i >= 0; i-- {
		mwOpt := r.middleware.Get(i)
		if mwOpt.IsSome() {
			mw := *mwOpt.Unwrap()
			capturedHandler := handler
			handler = func(ctx *Context) shepard.Result[Response, HttpError] {
				return mw(ctx, func() shepard.Result[Response, HttpError] {
					return capturedHandler(ctx)
				})
			}
		}
	}

	// Execute the handler chain
	result := handler(ctx)

	// Convert HttpError to Response if needed
	if result.IsErr() {
		err := result.UnwrapErr()
		return shepard.Ok[Response, HttpError](err.IntoResponse())
	}

	return result
}

// findHandler finds the matching handler for a request
func (r *Router) findHandler(ctx *Context) func(*Context) shepard.Result[Response, HttpError] {
	path := ctx.Path()
	method := ctx.Method()

	// Check if any route matches the path (regardless of method)
	var pathMatchFound bool
	var matchedRoute *Route

	// Find matching route
	routeIter := r.routes.Iter()
	for {
		routeOpt := routeIter.Next()
		if routeOpt.IsNone() {
			break
		}
		route := routeOpt.Unwrap()

		// Check if path matches (regardless of method)
		pathParams := route.Pattern.Match(path)
		if pathParams.IsSome() {
			pathMatchFound = true

			// Check if method also matches
			if route.Method == method {
				matchedRoute = &route

				// Copy params to context
				params := pathParams.Unwrap()
				params.KeysUnordered().Foreach(func(_ int, key string) {
					valOpt := params.Get(key)
					if valOpt.IsSome() {
						ctx.SetPathParam(key, *valOpt.Unwrap())
					}
				})

				break
			}
		}
	}

	if matchedRoute != nil {
		return func(ctx *Context) shepard.Result[Response, HttpError] {
			return matchedRoute.Handler.Handle(ctx)
		}
	}

	// Path matched but method didn't
	if pathMatchFound {
		return func(ctx *Context) shepard.Result[Response, HttpError] {
			return defaultMethodNotAllowedHandler().Handle(ctx)
		}
	}

	// No route found
	return func(ctx *Context) shepard.Result[Response, HttpError] {
		return r.notFoundHandler.Handle(ctx)
	}
}

// Middleware is a function that wraps a handler
type Middleware func(ctx *Context, next func() shepard.Result[Response, HttpError]) shepard.Result[Response, HttpError]

// --- Route Groups ---

// RouteGroup represents a group of routes with a shared prefix and optional middleware
type RouteGroup struct {
	router     *Router
	prefix     string
	middleware slice.Slice[Middleware]
}

// Group creates a new route group with the given prefix
// Routes added to the group will have the prefix prepended to their patterns
func (r *Router) Group(prefix string) *RouteGroup {
	return &RouteGroup{
		router:     r,
		prefix:     prefix,
		middleware: slice.New[Middleware](),
	}
}

// Use adds middleware to this group only
// Middleware is applied to all routes added to this group
func (g *RouteGroup) Use(mw Middleware) *RouteGroup {
	g.middleware.Push(mw)
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
	return &RouteGroup{
		router:     g.router,
		prefix:     g.prefix + prefix,
		middleware: g.middleware.Clone(),
	}
}

// addRoute adds a route to the parent router with the group's prefix and middleware
func (g *RouteGroup) addRoute(method Method, pattern string, handler Handler) {
	fullPattern := g.prefix + pattern

	// If no group middleware, add route directly
	if g.middleware.Len() == 0 {
		g.router.Route(method, fullPattern, handler)
		return
	}

	// Wrap handler with group middleware (innermost to outermost)
	wrappedHandler := handler
	for i := g.middleware.Len() - 1; i >= 0; i-- {
		mwOpt := g.middleware.Get(i)
		if mwOpt.IsSome() {
			mw := *mwOpt.Unwrap()
			capturedHandler := wrappedHandler
			wrappedHandler = HandlerFunc(func(ctx *Context) shepard.Result[Response, HttpError] {
				return mw(ctx, func() shepard.Result[Response, HttpError] {
					return capturedHandler.Handle(ctx)
				})
			})
		}
	}

	g.router.Route(method, fullPattern, wrappedHandler)
}
