package http

import (
	"net"

	"github.com/marlaone/shepard"
)

func HandleRequest(conn net.Conn, req Request[RequestBody]) shepard.Result[Response[Body], error] {

	// TODO add router

	// example response
	res := NewHttpResponseBytes()
	res.SetStatusCode(StatusCodeOk)
	res.SetHeader("Content-Type", "text/plain")
	res.Body().Write() <- []byte("Hello world!")
	res.Body().Finish()

	return shepard.Ok[Response[Body], error](res)
}
