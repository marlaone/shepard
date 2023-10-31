package http

type Handler func(req Request[RequestBody]) Response[Body]

type Route struct {
	Method  Method
	Pattern string
	Handler Handler
}

func Get(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodGet,
		Pattern: pattern,
		Handler: handler,
	}
}

func Post(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodPost,
		Pattern: pattern,
		Handler: handler,
	}
}

func Put(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodPut,
		Pattern: pattern,
		Handler: handler,
	}
}

func Delete(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodDelete,
		Pattern: pattern,
		Handler: handler,
	}
}

func Patch(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodPatch,
		Pattern: pattern,
		Handler: handler,
	}
}

func Head(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodHead,
		Pattern: pattern,
		Handler: handler,
	}
}

func Options(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodOptions,
		Pattern: pattern,
		Handler: handler,
	}
}

func Trace(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodTrace,
		Pattern: pattern,
		Handler: handler,
	}
}

func Connect(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodConnect,
		Pattern: pattern,
		Handler: handler,
	}
}
type Handler func(req Request[RequestBody]) Response[Body]

type Route struct {
	Method  Method
	Pattern string
	Handler Handler
}

func Get(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodGet,
		Pattern: pattern,
		Handler: handler,
	}
}

func Post(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodPost,
		Pattern: pattern,
		Handler: handler,
	}
}

func Put(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodPut,
		Pattern: pattern,
		Handler: handler,
	}
}

func Delete(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodDelete,
		Pattern: pattern,
		Handler: handler,
	}
}

func Patch(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodPatch,
		Pattern: pattern,
		Handler: handler,
	}
}

func Head(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodHead,
		Pattern: pattern,
		Handler: handler,
	}
}

func Options(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodOptions,
		Pattern: pattern,
		Handler: handler,
	}
}

func Trace(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodTrace,
		Pattern: pattern,
		Handler: handler,
	}
}

func Connect(pattern string, handler Handler) Route {
	return Route{
		Method:  MethodConnect,
		Pattern: pattern,
		Handler: handler,
	}
}
