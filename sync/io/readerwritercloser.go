package io

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/slice"
)

type Reader interface {
	Read(*slice.Slice[byte]) shepard.Result[int, error]
}

type Writer interface {
	Write(slice.Slice[byte]) shepard.Result[int, error]
}

type Closer interface {
	Close() shepard.Result[shepard.Nil, error]
}

type ReaderWriter interface {
	Reader
	Writer
}

type ReaderCloser interface {
	Reader
	Closer
}

type WriterCloser interface {
	Writer
	Closer
}

type ReaderWriterCloser interface {
	Reader
	Writer
	Closer
}
