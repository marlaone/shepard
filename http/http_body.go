package http

type BytesBody struct {
	c chan []byte
}

var _ Body = (*BytesBody)(nil)

func NewBytesBody() *BytesBody {
	return &BytesBody{
		c: make(chan []byte, 1),
	}
}

func (b *BytesBody) Finish() {
	close(b.c)
}

func (b *BytesBody) Read() <-chan []byte {
	return b.c
}

func (b *BytesBody) Write() chan<- []byte {
	return b.c
}
