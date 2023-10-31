package http

import "github.com/marlaone/shepard"

type Next func() shepard.Result[Response[Body], error]

type Middleware func(req *Request[RequestBody], next Next) shepard.Result[Response[Body], error]
