package io

import (
	"errors"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/slice"
	"github.com/marlaone/shepard/sync/rwmutex"
)

type ReaderWriterCloserBuffer struct {
	buf rwmutex.RWMutex[slice.Slice[byte]]

	closed bool
}

var _ ReaderWriterCloser = (*ReaderWriterCloserBuffer)(nil)

func NewReaderWriterCloser() *ReaderWriterCloserBuffer {
	return &ReaderWriterCloserBuffer{
		buf:    rwmutex.New[slice.Slice[byte]](slice.New[byte]()),
		closed: false,
	}
}

func (b *ReaderWriterCloserBuffer) Read(p *slice.Slice[byte]) shepard.Result[int, error] {
	guard := b.buf.RLock()
	defer b.buf.RUnlock()
	buf := guard.Unwrap()

	p.Append(buf)

	return shepard.Ok[int, error](p.Len())
}

func (b *ReaderWriterCloserBuffer) Write(p slice.Slice[byte]) shepard.Result[int, error] {

	if b.closed {
		return shepard.Err[int, error](errors.New("can't write to closed buffer"))
	}

	guard := b.buf.Lock()
	defer b.buf.Unlock()
	buf := guard.Unwrap()

	buf.Append(&p)

	return shepard.Ok[int, error](buf.Len())
}

func (b *ReaderWriterCloserBuffer) Close() shepard.Result[shepard.Nil, error] {
	b.closed = true
	return shepard.Ok[shepard.Nil, error](shepard.Nil{})
}

func (b *ReaderWriterCloserBuffer) Closed() bool {
	return b.closed
}
