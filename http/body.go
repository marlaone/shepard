package http

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/marlaone/shepard"
)

// BodyReader wraps an io.Reader for reading request bodies
// with optional size limit and Result-based error handling
type BodyReader struct {
	reader    io.Reader
	limit     shepard.Option[int64]
	bytesRead int64
}

// NewBodyReader creates a new BodyReader from an io.Reader
func NewBodyReader(r io.Reader) *BodyReader {
	return &BodyReader{
		reader:    r,
		limit:     shepard.None[int64](),
		bytesRead: 0,
	}
}

// NewBodyReaderWithLimit creates a BodyReader with a maximum read limit
func NewBodyReaderWithLimit(r io.Reader, limit int64) *BodyReader {
	return &BodyReader{
		reader:    io.LimitReader(r, limit),
		limit:     shepard.Some(limit),
		bytesRead: 0,
	}
}

// NewBodyReaderFromBytes creates a BodyReader from a byte slice
func NewBodyReaderFromBytes(data []byte) *BodyReader {
	return NewBodyReader(bytes.NewReader(data))
}

// EmptyBodyReader creates an empty BodyReader
func EmptyBodyReader() *BodyReader {
	return NewBodyReader(bytes.NewReader(nil))
}

// Read reads up to len(buf) bytes into buf
// Returns Result[int, error] with the number of bytes read
func (b *BodyReader) Read(buf []byte) shepard.Result[int, error] {
	n, err := b.reader.Read(buf)
	b.bytesRead += int64(n)
	if err != nil && err != io.EOF {
		return shepard.Err[int, error](err)
	}
	return shepard.Ok[int, error](n)
}

// ReadAll reads all remaining data from the body
// Returns Result[[]byte, error]
func (b *BodyReader) ReadAll() shepard.Result[[]byte, error] {
	data, err := io.ReadAll(b.reader)
	if err != nil {
		return shepard.Err[[]byte, error](err)
	}
	b.bytesRead += int64(len(data))
	return shepard.Ok[[]byte, error](data)
}

// BytesRead returns the number of bytes read so far
func (b *BodyReader) BytesRead() int64 {
	return b.bytesRead
}

// Limit returns the optional read limit
func (b *BodyReader) Limit() shepard.Option[int64] {
	return b.limit
}

// BodyJson deserializes the body as JSON into type T
// This is a free function to work around Go's generic method limitations
func BodyJson[T any](body *BodyReader) shepard.Result[T, error] {
	data := body.ReadAll()
	if data.IsErr() {
		return shepard.Err[T, error](data.UnwrapErr())
	}

	var value T
	if err := json.Unmarshal(data.Unwrap(), &value); err != nil {
		return shepard.Err[T, error](err)
	}

	return shepard.Ok[T, error](value)
}

// FormValues is a type alias for QueryValues (same structure: key -> []values)
// Used for parsing application/x-www-form-urlencoded request bodies
type FormValues = QueryValues

// BodyForm parses URL-encoded form data (application/x-www-form-urlencoded) from the body
// Returns Result[FormValues, error] - FormValues is the same structure as QueryValues
func BodyForm(body *BodyReader) shepard.Result[FormValues, error] {
	data := body.ReadAll()
	if data.IsErr() {
		return shepard.Err[FormValues, error](data.UnwrapErr())
	}

	// Reuse ParseQuery from uri.go - form encoding is the same as query string encoding
	return ParseQuery(string(data.Unwrap()))
}

// BodyWriter is a buffered writer for response bodies
// with Result-based error handling
type BodyWriter struct {
	buffer bytes.Buffer
}

// NewBodyWriter creates a new empty BodyWriter
func NewBodyWriter() *BodyWriter {
	return &BodyWriter{}
}

// Write writes data to the body buffer
// Returns Result[int, error] with the number of bytes written
func (b *BodyWriter) Write(data []byte) shepard.Result[int, error] {
	n, err := b.buffer.Write(data)
	if err != nil {
		return shepard.Err[int, error](err)
	}
	return shepard.Ok[int, error](n)
}

// WriteString writes a string to the body buffer
func (b *BodyWriter) WriteString(s string) shepard.Result[int, error] {
	n, err := b.buffer.WriteString(s)
	if err != nil {
		return shepard.Err[int, error](err)
	}
	return shepard.Ok[int, error](n)
}

// Json serializes a value as JSON and writes it to the body
func (b *BodyWriter) Json(value any) shepard.Result[shepard.Nil, error] {
	data, err := json.Marshal(value)
	if err != nil {
		return shepard.Err[shepard.Nil, error](err)
	}
	_, err = b.buffer.Write(data)
	if err != nil {
		return shepard.Err[shepard.Nil, error](err)
	}
	return shepard.Ok[shepard.Nil, error](shepard.Nil{})
}

// Bytes returns the written bytes
func (b *BodyWriter) Bytes() []byte {
	return b.buffer.Bytes()
}

// Len returns the number of bytes written
func (b *BodyWriter) Len() int {
	return b.buffer.Len()
}

// Reset clears the buffer
func (b *BodyWriter) Reset() {
	b.buffer.Reset()
}

// Reader returns an io.Reader for the written data
func (b *BodyWriter) Reader() io.Reader {
	return bytes.NewReader(b.buffer.Bytes())
}

// WriteTo writes the buffer to an io.Writer
func (b *BodyWriter) WriteTo(w io.Writer) (int64, error) {
	return b.buffer.WriteTo(w)
}
