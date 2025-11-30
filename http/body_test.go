package http

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// errorReader is a helper that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

// BodyReader tests

func TestBodyReader_Read(t *testing.T) {
	data := "Hello, World!"
	reader := NewBodyReader(strings.NewReader(data))

	buf := make([]byte, 5)
	result := reader.Read(buf)

	assert.True(t, result.IsOk())
	assert.Equal(t, 5, result.Unwrap())
	assert.Equal(t, "Hello", string(buf))
}

func TestBodyReader_ReadAll(t *testing.T) {
	data := "Hello, World!"
	reader := NewBodyReader(strings.NewReader(data))

	result := reader.ReadAll()

	assert.True(t, result.IsOk())
	assert.Equal(t, []byte(data), result.Unwrap())
}

func TestBodyReader_ReadAllWithLimit(t *testing.T) {
	data := "Hello, World! This is a longer message."
	reader := NewBodyReaderWithLimit(strings.NewReader(data), 5)

	result := reader.ReadAll()

	assert.True(t, result.IsOk())
	// Should only read up to the limit
	assert.Equal(t, 5, len(result.Unwrap()))
}

func TestBodyReader_Empty(t *testing.T) {
	reader := NewBodyReader(strings.NewReader(""))

	result := reader.ReadAll()

	assert.True(t, result.IsOk())
	assert.Equal(t, 0, len(result.Unwrap()))
}

func TestBodyJson(t *testing.T) {
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	jsonData := `{"name": "John", "email": "john@example.com"}`
	reader := NewBodyReader(strings.NewReader(jsonData))

	result := BodyJson[User](reader)

	assert.True(t, result.IsOk())
	user := result.Unwrap()
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
}

func TestBodyJson_InvalidJson(t *testing.T) {
	reader := NewBodyReader(strings.NewReader("not valid json"))

	result := BodyJson[map[string]string](reader)

	assert.True(t, result.IsErr())
}

// BodyWriter tests

func TestBodyWriter_Write(t *testing.T) {
	writer := NewBodyWriter()

	result := writer.Write([]byte("Hello"))

	assert.True(t, result.IsOk())
	assert.Equal(t, 5, result.Unwrap())
	assert.Equal(t, []byte("Hello"), writer.Bytes())
}

func TestBodyWriter_WriteString(t *testing.T) {
	writer := NewBodyWriter()

	result := writer.WriteString("Hello, World!")

	assert.True(t, result.IsOk())
	assert.Equal(t, 13, result.Unwrap())
	assert.Equal(t, "Hello, World!", string(writer.Bytes()))
}

func TestBodyWriter_WriteMultiple(t *testing.T) {
	writer := NewBodyWriter()

	writer.WriteString("Hello, ")
	writer.WriteString("World!")

	assert.Equal(t, "Hello, World!", string(writer.Bytes()))
}

func TestBodyWriter_Json(t *testing.T) {
	type Response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	writer := NewBodyWriter()

	result := writer.Json(Response{Message: "OK", Code: 200})

	assert.True(t, result.IsOk())

	expected := `{"message":"OK","code":200}`
	assert.Equal(t, expected, string(writer.Bytes()))
}

func TestBodyWriter_Json_Invalid(t *testing.T) {
	writer := NewBodyWriter()

	// Channels can't be marshaled to JSON
	ch := make(chan int)
	result := writer.Json(ch)

	assert.True(t, result.IsErr())
}

func TestBodyWriter_Len(t *testing.T) {
	writer := NewBodyWriter()
	assert.Equal(t, 0, writer.Len())

	writer.WriteString("Hello")
	assert.Equal(t, 5, writer.Len())
}

func TestBodyWriter_Reset(t *testing.T) {
	writer := NewBodyWriter()
	writer.WriteString("Hello")

	writer.Reset()

	assert.Equal(t, 0, writer.Len())
	assert.Equal(t, "", string(writer.Bytes()))
}

func TestBodyWriter_Reader(t *testing.T) {
	writer := NewBodyWriter()
	writer.WriteString("Hello, World!")

	reader := writer.Reader()

	data, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(data))
}

func TestEmptyBodyReader(t *testing.T) {
	reader := EmptyBodyReader()

	result := reader.ReadAll()

	assert.True(t, result.IsOk())
	assert.Equal(t, 0, len(result.Unwrap()))
}

func TestBodyReader_FromBytes(t *testing.T) {
	data := []byte("test data")
	reader := NewBodyReaderFromBytes(data)

	result := reader.ReadAll()

	assert.True(t, result.IsOk())
	assert.Equal(t, data, result.Unwrap())
}

func TestBodyWriter_WriteTo(t *testing.T) {
	writer := NewBodyWriter()
	writer.WriteString("Hello")

	var buf bytes.Buffer
	n, err := writer.WriteTo(&buf)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), n)
	assert.Equal(t, "Hello", buf.String())
}

// BodyForm tests

func TestBodyForm_Simple(t *testing.T) {
	formData := "name=John&email=john%40example.com"
	reader := NewBodyReader(strings.NewReader(formData))

	result := BodyForm(reader)

	assert.True(t, result.IsOk())
	form := result.Unwrap()

	nameOpt := form.Get("name")
	assert.True(t, nameOpt.IsSome())
	assert.Equal(t, "John", *nameOpt.Unwrap())

	emailOpt := form.Get("email")
	assert.True(t, emailOpt.IsSome())
	assert.Equal(t, "john@example.com", *emailOpt.Unwrap())
}

func TestBodyForm_MultipleValues(t *testing.T) {
	formData := "tag=go&tag=rust&tag=python"
	reader := NewBodyReader(strings.NewReader(formData))

	result := BodyForm(reader)

	assert.True(t, result.IsOk())
	form := result.Unwrap()

	// Get should return first value
	tagOpt := form.Get("tag")
	assert.True(t, tagOpt.IsSome())
	assert.Equal(t, "go", *tagOpt.Unwrap())

	// GetAll should return all values
	tags := form.GetAll("tag")
	assert.Equal(t, 3, tags.Len())
}

func TestBodyForm_Empty(t *testing.T) {
	reader := NewBodyReader(strings.NewReader(""))

	result := BodyForm(reader)

	assert.True(t, result.IsOk())
	form := result.Unwrap()
	assert.False(t, form.Has("anything"))
}

func TestBodyForm_SpecialChars(t *testing.T) {
	// URL-encoded special characters
	formData := "message=Hello+World&data=%26%3D%3F"
	reader := NewBodyReader(strings.NewReader(formData))

	result := BodyForm(reader)

	assert.True(t, result.IsOk())
	form := result.Unwrap()

	msgOpt := form.Get("message")
	assert.True(t, msgOpt.IsSome())
	assert.Equal(t, "Hello World", *msgOpt.Unwrap())

	dataOpt := form.Get("data")
	assert.True(t, dataOpt.IsSome())
	assert.Equal(t, "&=?", *dataOpt.Unwrap())
}

func TestBodyForm_ReadError(t *testing.T) {
	// Create a reader that will fail
	reader := NewBodyReader(&errorReader{})

	result := BodyForm(reader)

	assert.True(t, result.IsErr())
}
