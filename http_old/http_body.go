package http

import (
	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/slice"
	"github.com/marlaone/shepard/sync/io"
)

type BytesBody struct {
	rwc *io.ReaderWriterCloserBuffer
}

var _ Body = (*BytesBody)(nil)

func NewBytesBody() *BytesBody {
	return &BytesBody{
		rwc: io.NewReaderWriterCloser(),
	}
}

func (b *BytesBody) Close() shepard.Result[shepard.Nil, error] {
	return b.rwc.Close()
}

func (b *BytesBody) Read(bs *slice.Slice[byte]) shepard.Result[int, error] {
	return b.rwc.Read(bs)
}

func (b *BytesBody) Write(bs slice.Slice[byte]) shepard.Result[int, error] {
	return b.rwc.Write(bs)
}

func (b *BytesBody) Closed() bool {
	return b.rwc.Closed()
}
