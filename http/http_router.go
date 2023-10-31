package http

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/slice"
)

type Router struct {
	routes slice.Slice[Route]

	middleware slice.Slice[Middleware]

	notFoundHandler Handler
}

func NewRouter() *Router {
	return &Router{
		routes:     slice.New[Route](),
		middleware: slice.New[Middleware](),
		notFoundHandler: func(req *Request[RequestBody]) Response[Body] {
			res := NewResponseBuilder(NewHttpResponseBytes()).Status(404).Body(NewBytesBody()).Unwrap()
			res.Body().Finish()
			return res
		},
	}
}

func (r *Router) Use(middleware Middleware) {
	r.middleware.Push(middleware)
}

func (r *Router) Register(route Route) {
	r.routes.Push(route)
}

func (r *Router) NotFound(handler Handler) {
	r.notFoundHandler = handler
}

func (r *Router) Serve(req *Request[RequestBody]) shepard.Result[Response[Body], error] {
	return invokeMiddlewares(r.middleware, req, func() shepard.Result[Response[Body], error] {
		route := r.routes.Iter().Find(func(val *Route) bool {
			// TODO find route by pattern
			return val.Pattern == req.URL.Path
		})
		if route.IsNone() {
			return shepard.Ok[Response[Body], error](r.notFoundHandler(req))
		}
		if route.Unwrap().Method != req.Method {
			res := NewResponseBuilder(NewHttpResponseBytes()).Status(405).Body(NewBytesBody()).Unwrap()
			res.Body().Finish()
			return shepard.Ok[Response[Body], error](res)
		}
		return shepard.Ok[Response[Body], error](route.Unwrap().Handler(req))
	})
}

func invokeMiddlewares(middleware slice.Slice[Middleware], req *Request[RequestBody], next Next) shepard.Result[Response[Body], error] {
	mwIter := middleware.Iter()
	nextMW := mwIter.Next()
	for nextMW.IsSome() {
		mw := nextMW.Unwrap()
		nextMW = mwIter.Next()
		res := mw(req, next)
		if res.IsErr() {
			return res
		}
	}
	return next()
}
