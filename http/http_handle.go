package http

import (
	"net"

	"github.com/marlaone/shepard"
)

func HandleRequest(conn net.Conn, req Request[RequestBody], r *Router) shepard.Result[Response[Body], error] {
	return r.Serve(&req)
}
