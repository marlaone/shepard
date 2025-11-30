package http

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
)

// DefaultMaxMemory is the default maximum memory for multipart parsing (32 MB)
const DefaultMaxMemory int64 = 32 << 20

// FilePart represents an uploaded file from a multipart form
type FilePart struct {
	Filename    string
	ContentType string
	Size        int64
	data        []byte
}

// Bytes returns the file content as a byte slice
func (f *FilePart) Bytes() []byte {
	return f.data
}

// Reader returns an io.Reader for the file content
func (f *FilePart) Reader() io.Reader {
	return bytes.NewReader(f.data)
}

// MultipartForm represents parsed multipart form data
type MultipartForm struct {
	Values FormValues
	Files  hashmap.HashMap[string, slice.Slice[FilePart]]
}

// Get returns the first file for the given field name
func (m *MultipartForm) Get(name string) shepard.Option[*FilePart] {
	filesOpt := m.Files.Get(name)
	if filesOpt.IsNone() {
		return shepard.None[*FilePart]()
	}
	files := filesOpt.Unwrap()
	firstOpt := files.First()
	if firstOpt.IsNone() {
		return shepard.None[*FilePart]()
	}
	return shepard.Some(firstOpt.Unwrap())
}

// GetAll returns all files for the given field name
func (m *MultipartForm) GetAll(name string) slice.Slice[FilePart] {
	filesOpt := m.Files.Get(name)
	if filesOpt.IsNone() {
		return slice.New[FilePart]()
	}
	return *filesOpt.Unwrap()
}

// Value returns the first form value for the given key
func (m *MultipartForm) Value(name string) shepard.Option[*string] {
	return m.Values.Get(name)
}

// BodyMultipart parses multipart form data from the body
// Uses DefaultMaxMemory (32 MB) as the maximum memory limit
func BodyMultipart(body *BodyReader, contentType string) shepard.Result[MultipartForm, error] {
	return BodyMultipartWithLimit(body, contentType, DefaultMaxMemory)
}

// BodyMultipartWithLimit parses multipart form data with a custom memory limit
// Files larger than maxMemory will be stored in temporary files
func BodyMultipartWithLimit(body *BodyReader, contentType string, maxMemory int64) shepard.Result[MultipartForm, error] {
	// Parse the Content-Type header to get the boundary
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return shepard.Err[MultipartForm, error](err)
	}

	if mediaType != "multipart/form-data" {
		return shepard.Err[MultipartForm, error](
			&multipartError{msg: "content type is not multipart/form-data"},
		)
	}

	boundary, ok := params["boundary"]
	if !ok {
		return shepard.Err[MultipartForm, error](
			&multipartError{msg: "missing boundary parameter"},
		)
	}

	// Create multipart reader
	mr := multipart.NewReader(body.reader, boundary)

	// Parse the form
	form, err := mr.ReadForm(maxMemory)
	if err != nil {
		return shepard.Err[MultipartForm, error](err)
	}

	// Convert to our types
	result := MultipartForm{
		Values: NewQueryValues(),
		Files:  hashmap.New[string, slice.Slice[FilePart]](),
	}

	// Copy form values
	for key, values := range form.Value {
		for _, v := range values {
			result.Values.Add(key, v)
		}
	}

	// Copy files
	for key, fileHeaders := range form.File {
		files := slice.New[FilePart]()
		for _, fh := range fileHeaders {
			// Read file content
			f, err := fh.Open()
			if err != nil {
				// Clean up the form and return error
				form.RemoveAll()
				return shepard.Err[MultipartForm, error](err)
			}
			data, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				form.RemoveAll()
				return shepard.Err[MultipartForm, error](err)
			}

			filePart := FilePart{
				Filename:    fh.Filename,
				ContentType: fh.Header.Get("Content-Type"),
				Size:        int64(len(data)),
				data:        data,
			}
			files.Push(filePart)
		}
		result.Files.Insert(key, files)
	}

	// Clean up temporary files created by multipart.ReadForm
	form.RemoveAll()

	return shepard.Ok[MultipartForm, error](result)
}

// multipartError is a custom error type for multipart parsing errors
type multipartError struct {
	msg string
}

func (e *multipartError) Error() string {
	return e.msg
}
