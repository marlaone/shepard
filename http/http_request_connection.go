package http

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/marlaone/shepard"
)

var headerKeySeparator = []byte{':'}
var headerValueSeparator = []byte{','}

func RequestFromConnection(conn net.Conn) shepard.Result[Request[RequestBody], error] {

	// create buffer
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return shepard.Err[Request[RequestBody], error](fmt.Errorf("[http.RequestFromConnection] read failed: %w", err))
	}

	// check if creating buffer was successful
	if n == 0 {
		return shepard.Ok[Request[RequestBody], error](Request[RequestBody]{}.Default())
	}

	// start parsing

	buffer := bytes.NewBuffer(buf)

	// read method
	m, err := buffer.ReadBytes(' ')
	if err != nil {
		return shepard.Err[Request[RequestBody], error](fmt.Errorf("[http.RequestFromConnection] read method failed: %w", err))
	}

	// check if method is valid
	method := TryMethodFromString(string(m))
	if method.IsErr() {
		return shepard.Err[Request[RequestBody], error](fmt.Errorf("[http.RequestFromConnection] invalid method: %w", method.Err().Unwrap()))
	}

	// read path
	// TODO: parse path
	_, err = buffer.ReadBytes(' ')
	if err != nil {
		return shepard.Err[Request[RequestBody], error](fmt.Errorf("[http.RequestFromConnection] read path failed: %w", err))
	}

	// read protocol
	protocol, err := buffer.ReadBytes('\n')
	if err != nil {
		return shepard.Err[Request[RequestBody], error](fmt.Errorf("[http.RequestFromConnection] read protocol failed: %w", err))
	}
	protocol = bytes.TrimSpace(protocol)

	// create request builder
	builder := NewRequestBuilder[RequestBody]().Method(method.Unwrap()).Version(Version(protocol))

	// read headers
	for {
		line, err := buffer.ReadBytes('\n')
		if err != nil {
			return shepard.Err[Request[RequestBody], error](fmt.Errorf("[http.RequestFromConnection] read header line failed: %w", err))
		}

		// check if line is empty
		if len(line) == 1 {
			break
		}

		// parse header
		linesBuffer := bytes.NewBuffer(line)

		// read header key
		key, err := linesBuffer.ReadBytes(':')
		if err != nil {
			if err == io.EOF {
				break
			}
			return shepard.Err[Request[RequestBody], error](fmt.Errorf("[http.RequestFromConnection] read header key failed: %w", err))
		}

		// read header value
		headerValue, err := linesBuffer.ReadBytes('\n')
		if err != nil {
			return shepard.Err[Request[RequestBody], error](fmt.Errorf("[http.RequestFromConnection] read header values failed: %w", err))
		}

		// parse header value to slice
		splitted := bytes.Split(headerValue, headerValueSeparator)
		values := make([]string, 0, len(splitted))
		for _, value := range splitted {
			values = append(values, string(bytes.TrimSpace(value)))
		}

		// add header to request
		builder.Header(string(bytes.TrimSpace(key)), values...)
	}

	// add remaining bytes from buffer to request body and return request
	return builder.Body(buffer)
}
