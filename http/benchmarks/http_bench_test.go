// Package benchmarks provides comprehensive benchmarks for the HTTP package
package benchmarks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/http"
	"github.com/marlaone/shepard/http/extract"
)

// ============================================================================
// Test Data Types
// ============================================================================

type testUser struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Age      int      `json:"age"`
	Active   bool     `json:"active"`
	Tags     []string `json:"tags"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ============================================================================
// Helper Functions
// ============================================================================

// makeHeaders creates n custom headers
func makeHeaders(n int) *http.Headers {
	h := http.NewHeaders()
	for i := 0; i < n; i++ {
		h.Set(fmt.Sprintf("X-Custom-Header-%d", i), fmt.Sprintf("value-%d", i))
	}
	return h
}

// makeQueryString creates query string with n params
func makeQueryString(n int) string {
	params := make([]string, n)
	for i := 0; i < n; i++ {
		params[i] = fmt.Sprintf("param%d=value%d", i, i)
	}
	return strings.Join(params, "&")
}

// makeTestUser creates a single test user
func makeTestUser(id int) testUser {
	return testUser{
		ID:     id,
		Name:   fmt.Sprintf("User %d", id),
		Email:  fmt.Sprintf("user%d@example.com", id),
		Age:    20 + (id % 50),
		Active: id%2 == 0,
		Tags:   []string{"tag1", "tag2", "tag3"},
	}
}

// makeTestUsers creates n test users
func makeTestUsers(n int) []testUser {
	users := make([]testUser, n)
	for i := 0; i < n; i++ {
		users[i] = makeTestUser(i)
	}
	return users
}

// makeJSONBytes creates JSON bytes for n users
func makeJSONBytes(n int) []byte {
	users := makeTestUsers(n)
	data, _ := json.Marshal(users)
	return data
}

// makeSingleUserJSON creates JSON bytes for a single user
func makeSingleUserJSON() []byte {
	user := makeTestUser(1)
	data, _ := json.Marshal(user)
	return data
}

// makeRouterWithNRoutes creates router with n static routes
func makeRouterWithNRoutes(n int) *http.Router {
	r := http.NewRouter()
	for i := 0; i < n; i++ {
		path := fmt.Sprintf("/api/v1/resource%d", i)
		r.Get(path, http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
			return http.Text("ok")
		}))
	}
	return r
}

// makeFormData creates URL-encoded form data
func makeFormData(numFields int) string {
	fields := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		fields[i] = fmt.Sprintf("field%d=value%d", i, i)
	}
	return strings.Join(fields, "&")
}

// makeFormDataWithLongValues creates form data with long values
func makeFormDataWithLongValues(numFields int, valueLen int) string {
	fields := make([]string, numFields)
	longValue := strings.Repeat("x", valueLen)
	for i := 0; i < numFields; i++ {
		fields[i] = fmt.Sprintf("field%d=%s%d", i, longValue, i)
	}
	return strings.Join(fields, "&")
}

// makeMultipartBody creates a multipart form body with files
func makeMultipartBody(fileSize int, numFiles int) (*bytes.Buffer, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	fileContent := make([]byte, fileSize)
	for i := range fileContent {
		fileContent[i] = byte(i % 256)
	}

	for i := 0; i < numFiles; i++ {
		part, _ := writer.CreateFormFile("file", fmt.Sprintf("file%d.bin", i))
		part.Write(fileContent)
	}

	writer.Close()
	return body, writer.FormDataContentType()
}

// makeMultipartBodyMixed creates a multipart form with both files and fields
func makeMultipartBodyMixed(fileSize int, numFiles int, numFields int) (*bytes.Buffer, string) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add fields
	for i := 0; i < numFields; i++ {
		writer.WriteField(fmt.Sprintf("field%d", i), fmt.Sprintf("value%d", i))
	}

	// Add files
	fileContent := make([]byte, fileSize)
	for i := range fileContent {
		fileContent[i] = byte(i % 256)
	}

	for i := 0; i < numFiles; i++ {
		part, _ := writer.CreateFormFile("file", fmt.Sprintf("file%d.bin", i))
		part.Write(fileContent)
	}

	writer.Close()
	return body, writer.FormDataContentType()
}

// ============================================================================
// Request Building Benchmarks
// ============================================================================

func BenchmarkRequest_Build_Simple(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	}
}

func BenchmarkRequest_Build_WithPath(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = http.NewRequest(http.MethodGet, "/api/v1/users/123/posts", http.EmptyBodyReader())
	}
}

func BenchmarkRequest_Build_WithQuery(b *testing.B) {
	path := "/api/v1/users?" + makeQueryString(10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = http.NewRequest(http.MethodGet, path, http.EmptyBodyReader())
	}
}

func BenchmarkRequest_Builder(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = http.NewRequestBuilder().
			Method(http.MethodPost).
			Path("/api/v1/users").
			Header("Content-Type", "application/json").
			Header("Authorization", "Bearer token123").
			Build()
	}
}

// ============================================================================
// Response Building Benchmarks
// ============================================================================

func BenchmarkResponse_NewResponse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = http.NewResponse(http.StatusOK)
	}
}

func BenchmarkResponse_Text_100B(b *testing.B) {
	text := strings.Repeat("x", 100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = http.Text(text)
	}
}

func BenchmarkResponse_Text_1KB(b *testing.B) {
	text := strings.Repeat("x", 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = http.Text(text)
	}
}

func BenchmarkResponse_Text_10KB(b *testing.B) {
	text := strings.Repeat("x", 10*1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = http.Text(text)
	}
}

func BenchmarkResponse_Builder(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = http.NewResponseBuilder().
			Status(http.StatusOK).
			Header("Content-Type", "text/plain").
			Text("Hello, World!").
			Build()
	}
}

func BenchmarkResponse_Builder_WithHeaders_5(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := http.NewResponseBuilder().Status(http.StatusOK)
		for j := 0; j < 5; j++ {
			builder.Header(fmt.Sprintf("X-Header-%d", j), fmt.Sprintf("value-%d", j))
		}
		_ = builder.Build()
	}
}

func BenchmarkResponse_Builder_WithHeaders_20(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := http.NewResponseBuilder().Status(http.StatusOK)
		for j := 0; j < 20; j++ {
			builder.Header(fmt.Sprintf("X-Header-%d", j), fmt.Sprintf("value-%d", j))
		}
		_ = builder.Build()
	}
}

// ============================================================================
// JSON Benchmarks
// ============================================================================

func BenchmarkResponse_Json_Small(b *testing.B) {
	user := makeTestUser(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := http.Json(user)
		if result.IsErr() {
			b.Fatal("json failed")
		}
	}
}

func BenchmarkResponse_Json_Medium(b *testing.B) {
	users := makeTestUsers(10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := http.Json(users)
		if result.IsErr() {
			b.Fatal("json failed")
		}
	}
}

func BenchmarkResponse_Json_Large(b *testing.B) {
	users := makeTestUsers(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := http.Json(users)
		if result.IsErr() {
			b.Fatal("json failed")
		}
	}
}

func BenchmarkBodyJson_Small(b *testing.B) {
	jsonData := makeSingleUserJSON()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(jsonData)
		result := http.BodyJson[testUser](reader)
		if result.IsErr() {
			b.Fatal("json parse failed")
		}
	}
}

func BenchmarkBodyJson_Medium(b *testing.B) {
	jsonData := makeJSONBytes(10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(jsonData)
		result := http.BodyJson[[]testUser](reader)
		if result.IsErr() {
			b.Fatal("json parse failed")
		}
	}
}

func BenchmarkBodyJson_Large(b *testing.B) {
	jsonData := makeJSONBytes(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(jsonData)
		result := http.BodyJson[[]testUser](reader)
		if result.IsErr() {
			b.Fatal("json parse failed")
		}
	}
}

func BenchmarkIntoResponse_JsonBody(b *testing.B) {
	user := makeTestUser(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jr := http.JsonBody(user)
		_ = jr.IntoResponse()
	}
}

// ============================================================================
// Form Parsing Benchmarks
// ============================================================================

func BenchmarkBodyForm_5Fields(b *testing.B) {
	formData := makeFormData(5)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes([]byte(formData))
		result := http.BodyForm(reader)
		if result.IsErr() {
			b.Fatal("form parse failed")
		}
	}
}

func BenchmarkBodyForm_20Fields(b *testing.B) {
	formData := makeFormData(20)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes([]byte(formData))
		result := http.BodyForm(reader)
		if result.IsErr() {
			b.Fatal("form parse failed")
		}
	}
}

func BenchmarkBodyForm_50Fields(b *testing.B) {
	formData := makeFormData(50)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes([]byte(formData))
		result := http.BodyForm(reader)
		if result.IsErr() {
			b.Fatal("form parse failed")
		}
	}
}

func BenchmarkBodyForm_LongValues(b *testing.B) {
	formData := makeFormDataWithLongValues(10, 1000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes([]byte(formData))
		result := http.BodyForm(reader)
		if result.IsErr() {
			b.Fatal("form parse failed")
		}
	}
}

// ============================================================================
// Multipart Parsing Benchmarks
// ============================================================================

func BenchmarkMultipart_1File_1KB(b *testing.B) {
	body, contentType := makeMultipartBody(1024, 1)
	bodyBytes := body.Bytes()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(bodyBytes)
		result := http.BodyMultipart(reader, contentType)
		if result.IsErr() {
			b.Fatal("multipart parse failed")
		}
	}
}

func BenchmarkMultipart_1File_100KB(b *testing.B) {
	body, contentType := makeMultipartBody(100*1024, 1)
	bodyBytes := body.Bytes()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(bodyBytes)
		result := http.BodyMultipart(reader, contentType)
		if result.IsErr() {
			b.Fatal("multipart parse failed")
		}
	}
}

func BenchmarkMultipart_1File_1MB(b *testing.B) {
	body, contentType := makeMultipartBody(1024*1024, 1)
	bodyBytes := body.Bytes()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(bodyBytes)
		result := http.BodyMultipart(reader, contentType)
		if result.IsErr() {
			b.Fatal("multipart parse failed")
		}
	}
}

func BenchmarkMultipart_5Files_10KB(b *testing.B) {
	body, contentType := makeMultipartBody(10*1024, 5)
	bodyBytes := body.Bytes()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(bodyBytes)
		result := http.BodyMultipart(reader, contentType)
		if result.IsErr() {
			b.Fatal("multipart parse failed")
		}
	}
}

func BenchmarkMultipart_Mixed(b *testing.B) {
	body, contentType := makeMultipartBodyMixed(10*1024, 2, 10)
	bodyBytes := body.Bytes()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(bodyBytes)
		result := http.BodyMultipart(reader, contentType)
		if result.IsErr() {
			b.Fatal("multipart parse failed")
		}
	}
}

// ============================================================================
// File Download Benchmarks
// ============================================================================

func BenchmarkFileResponse_Small(b *testing.B) {
	data := make([]byte, 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fr := http.File(data, "file.bin")
		_ = fr.IntoResponse()
	}
}

func BenchmarkFileResponse_Large(b *testing.B) {
	data := make([]byte, 1024*1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fr := http.File(data, "file.bin")
		_ = fr.IntoResponse()
	}
}

func BenchmarkFileResponse_WithContentType(b *testing.B) {
	data := make([]byte, 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fr := http.FileWithContentType(data, "document.pdf", "application/pdf")
		_ = fr.IntoResponse()
	}
}

// ============================================================================
// Router Benchmarks
// ============================================================================

func BenchmarkRouter_Match_Static_5Routes(b *testing.B) {
	router := makeRouterWithNRoutes(5)
	req := http.NewRequest(http.MethodGet, "/api/v1/resource2", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = router.Handle(ctx)
	}
}

func BenchmarkRouter_Match_Static_50Routes(b *testing.B) {
	router := makeRouterWithNRoutes(50)
	req := http.NewRequest(http.MethodGet, "/api/v1/resource25", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = router.Handle(ctx)
	}
}

func BenchmarkRouter_Match_Static_200Routes(b *testing.B) {
	router := makeRouterWithNRoutes(200)
	req := http.NewRequest(http.MethodGet, "/api/v1/resource100", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = router.Handle(ctx)
	}
}

func BenchmarkRouter_Match_WithParams_1Param(b *testing.B) {
	router := http.NewRouter()
	router.Get("/users/:id", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))
	req := http.NewRequest(http.MethodGet, "/users/123", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.PathParams.Clear()
		_ = router.Handle(ctx)
	}
}

func BenchmarkRouter_Match_WithParams_3Params(b *testing.B) {
	router := http.NewRouter()
	router.Get("/users/:userId/posts/:postId/comments/:commentId", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))
	req := http.NewRequest(http.MethodGet, "/users/123/posts/456/comments/789", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.PathParams.Clear()
		_ = router.Handle(ctx)
	}
}

func BenchmarkRouter_NotFound(b *testing.B) {
	router := makeRouterWithNRoutes(50)
	req := http.NewRequest(http.MethodGet, "/nonexistent/path", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = router.Handle(ctx)
	}
}

func BenchmarkRouter_MethodNotAllowed(b *testing.B) {
	router := http.NewRouter()
	router.Get("/users", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))
	req := http.NewRequest(http.MethodPost, "/users", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = router.Handle(ctx)
	}
}

// ============================================================================
// Handler Benchmarks
// ============================================================================

func BenchmarkHandler_Simple(b *testing.B) {
	handler := http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("Hello, World!")
	})
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(ctx)
	}
}

func BenchmarkHandler_WithResult(b *testing.B) {
	handler := http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		return shepard.Ok[http.Response, http.HttpError](http.Text("Hello, World!"))
	})
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(ctx)
	}
}

func BenchmarkHandler_WriterFunc(b *testing.B) {
	handler := http.WriterHandlerFunc(func(ctx *http.Context, w *http.ResponseWriter) {
		w.WriteString("Hello, World!")
	})
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.Handle(ctx)
	}
}

// ============================================================================
// Middleware Benchmarks
// ============================================================================

func BenchmarkHandler_WithMiddleware_1(b *testing.B) {
	router := http.NewRouter()
	router.Use(func(ctx *http.Context, next func() shepard.Result[http.Response, http.HttpError]) shepard.Result[http.Response, http.HttpError] {
		return next()
	})
	router.Get("/", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = router.Handle(ctx)
	}
}

func BenchmarkHandler_WithMiddleware_5(b *testing.B) {
	router := http.NewRouter()
	for j := 0; j < 5; j++ {
		router.Use(func(ctx *http.Context, next func() shepard.Result[http.Response, http.HttpError]) shepard.Result[http.Response, http.HttpError] {
			return next()
		})
	}
	router.Get("/", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = router.Handle(ctx)
	}
}

func BenchmarkHandler_WithMiddleware_10(b *testing.B) {
	router := http.NewRouter()
	for j := 0; j < 10; j++ {
		router.Use(func(ctx *http.Context, next func() shepard.Result[http.Response, http.HttpError]) shepard.Result[http.Response, http.HttpError] {
			return next()
		})
	}
	router.Get("/", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = router.Handle(ctx)
	}
}

// ============================================================================
// Extractor Benchmarks
// ============================================================================

func BenchmarkExtract_PathString(b *testing.B) {
	router := http.NewRouter()
	router.Get("/users/:id", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))
	req := http.NewRequest(http.MethodGet, "/users/123", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	router.Handle(ctx) // Populate path params

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := extract.ExtractPath[string](ctx, "id")
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

func BenchmarkExtract_PathInt(b *testing.B) {
	router := http.NewRouter()
	router.Get("/users/:id", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))
	req := http.NewRequest(http.MethodGet, "/users/123", http.EmptyBodyReader())
	ctx := http.NewContext(req)
	router.Handle(ctx) // Populate path params

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := extract.ExtractPath[int](ctx, "id")
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

func BenchmarkExtract_QueryString(b *testing.B) {
	req := http.NewRequest(http.MethodGet, "/users?name=John&age=30", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := extract.ExtractQuery[string](ctx, "name")
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

func BenchmarkExtract_QueryInt(b *testing.B) {
	req := http.NewRequest(http.MethodGet, "/users?name=John&age=30", http.EmptyBodyReader())
	ctx := http.NewContext(req)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := extract.ExtractQuery[int](ctx, "age")
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

func BenchmarkExtract_Header(b *testing.B) {
	req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader()).
		WithHeader("Authorization", "Bearer token123")
	ctx := http.NewContext(req)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := extract.ExtractHeader(ctx, "Authorization")
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

func BenchmarkExtract_Json_Small(b *testing.B) {
	jsonData := makeSingleUserJSON()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(jsonData)
		req := http.NewRequest(http.MethodPost, "/", reader).
			WithHeader("Content-Type", "application/json")
		ctx := http.NewContext(req)
		result := extract.Json[testUser](ctx)
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

func BenchmarkExtract_Json_Large(b *testing.B) {
	jsonData := makeJSONBytes(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(jsonData)
		req := http.NewRequest(http.MethodPost, "/", reader).
			WithHeader("Content-Type", "application/json")
		ctx := http.NewContext(req)
		result := extract.Json[[]testUser](ctx)
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

func BenchmarkExtract_Form(b *testing.B) {
	formData := makeFormData(10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes([]byte(formData))
		req := http.NewRequest(http.MethodPost, "/", reader).
			WithHeader("Content-Type", "application/x-www-form-urlencoded")
		ctx := http.NewContext(req)
		result := extract.Form(ctx)
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

func BenchmarkExtract_Multipart(b *testing.B) {
	body, contentType := makeMultipartBody(10*1024, 1)
	bodyBytes := body.Bytes()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(bodyBytes)
		req := http.NewRequest(http.MethodPost, "/", reader).
			WithHeader("Content-Type", contentType)
		ctx := http.NewContext(req)
		result := extract.Multipart(ctx)
		if result.IsErr() {
			b.Fatal("extract failed")
		}
	}
}

// ============================================================================
// Full Request/Response Cycle Benchmarks
// ============================================================================

func BenchmarkFullCycle_SimpleGet(b *testing.B) {
	router := http.NewRouter()
	router.Get("/", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("Hello, World!")
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader())
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}

func BenchmarkFullCycle_GetWithParams(b *testing.B) {
	router := http.NewRouter()
	router.Get("/users/:id", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		idOpt := ctx.Param("id")
		if idOpt.IsNone() {
			return http.TextWithStatus(http.StatusBadRequest, "missing id")
		}
		return http.Text("User: " + idOpt.Unwrap())
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := http.NewRequest(http.MethodGet, "/users/123", http.EmptyBodyReader())
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}

func BenchmarkFullCycle_PostJSON_Small(b *testing.B) {
	router := http.NewRouter()
	router.Post("/users", http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		result := extract.Json[testUser](ctx)
		if result.IsErr() {
			return shepard.Err[http.Response, http.HttpError](result.UnwrapErr())
		}
		jsonResult := http.Json(result.Unwrap())
		if jsonResult.IsErr() {
			return shepard.Err[http.Response, http.HttpError](http.InternalError(jsonResult.UnwrapErr()))
		}
		return shepard.Ok[http.Response, http.HttpError](jsonResult.Unwrap())
	}))

	jsonData := makeSingleUserJSON()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(jsonData)
		req := http.NewRequest(http.MethodPost, "/users", reader).
			WithHeader("Content-Type", "application/json")
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}

func BenchmarkFullCycle_PostJSON_Large(b *testing.B) {
	router := http.NewRouter()
	router.Post("/users", http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		result := extract.Json[[]testUser](ctx)
		if result.IsErr() {
			return shepard.Err[http.Response, http.HttpError](result.UnwrapErr())
		}
		jsonResult := http.Json(result.Unwrap())
		if jsonResult.IsErr() {
			return shepard.Err[http.Response, http.HttpError](http.InternalError(jsonResult.UnwrapErr()))
		}
		return shepard.Ok[http.Response, http.HttpError](jsonResult.Unwrap())
	}))

	jsonData := makeJSONBytes(100)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(jsonData)
		req := http.NewRequest(http.MethodPost, "/users", reader).
			WithHeader("Content-Type", "application/json")
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}

func BenchmarkFullCycle_WithMiddleware(b *testing.B) {
	router := http.NewRouter()
	// Add logging-like middleware
	router.Use(func(ctx *http.Context, next func() shepard.Result[http.Response, http.HttpError]) shepard.Result[http.Response, http.HttpError] {
		// Simulate logging setup
		_ = ctx.Method()
		_ = ctx.Path()
		return next()
	})
	// Add auth-like middleware
	router.Use(func(ctx *http.Context, next func() shepard.Result[http.Response, http.HttpError]) shepard.Result[http.Response, http.HttpError] {
		// Simulate auth check
		_ = ctx.Header("Authorization")
		return next()
	})
	router.Get("/", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("ok")
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := http.NewRequest(http.MethodGet, "/", http.EmptyBodyReader()).
			WithHeader("Authorization", "Bearer token123")
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}

func BenchmarkFullCycle_FormSubmit(b *testing.B) {
	router := http.NewRouter()
	router.Post("/submit", http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		result := extract.Form(ctx)
		if result.IsErr() {
			return shepard.Err[http.Response, http.HttpError](result.UnwrapErr())
		}
		form := result.Unwrap()
		nameOpt := form.Get("name")
		if nameOpt.IsNone() {
			return shepard.Err[http.Response, http.HttpError](http.BadRequest("missing name"))
		}
		return shepard.Ok[http.Response, http.HttpError](http.Text("Hello, " + *nameOpt.Unwrap()))
	}))

	formData := "name=John&email=john%40example.com&age=30"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes([]byte(formData))
		req := http.NewRequest(http.MethodPost, "/submit", reader).
			WithHeader("Content-Type", "application/x-www-form-urlencoded")
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}

func BenchmarkFullCycle_FileUpload(b *testing.B) {
	router := http.NewRouter()
	router.Post("/upload", http.HandlerFunc(func(ctx *http.Context) shepard.Result[http.Response, http.HttpError] {
		result := extract.Multipart(ctx)
		if result.IsErr() {
			return shepard.Err[http.Response, http.HttpError](result.UnwrapErr())
		}
		form := result.Unwrap()
		fileOpt := form.Get("file")
		if fileOpt.IsNone() {
			return shepard.Err[http.Response, http.HttpError](http.BadRequest("missing file"))
		}
		file := fileOpt.Unwrap()
		return shepard.Ok[http.Response, http.HttpError](
			http.Text(fmt.Sprintf("Uploaded: %s (%d bytes)", file.Filename, file.Size)),
		)
	}))

	body, contentType := makeMultipartBody(10*1024, 1) // 10KB file
	bodyBytes := body.Bytes()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(bodyBytes)
		req := http.NewRequest(http.MethodPost, "/upload", reader).
			WithHeader("Content-Type", contentType)
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}

func BenchmarkFullCycle_FileDownload(b *testing.B) {
	fileData := make([]byte, 100*1024) // 100KB file
	for i := range fileData {
		fileData[i] = byte(i % 256)
	}

	router := http.NewRouter()
	router.Get("/download", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.File(fileData, "data.bin").IntoResponse()
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := http.NewRequest(http.MethodGet, "/download", http.EmptyBodyReader())
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}

// ============================================================================
// URL/Query Parsing Benchmarks
// ============================================================================

func BenchmarkParseURL_Simple(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result := http.ParseURL("/api/v1/users")
		if result.IsErr() {
			b.Fatal("parse failed")
		}
	}
}

func BenchmarkParseURL_WithQuery(b *testing.B) {
	url := "/api/v1/users?" + makeQueryString(10)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := http.ParseURL(url)
		if result.IsErr() {
			b.Fatal("parse failed")
		}
	}
}

func BenchmarkParseQuery_5Params(b *testing.B) {
	query := makeQueryString(5)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := http.ParseQuery(query)
		if result.IsErr() {
			b.Fatal("parse failed")
		}
	}
}

func BenchmarkParseQuery_20Params(b *testing.B) {
	query := makeQueryString(20)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := http.ParseQuery(query)
		if result.IsErr() {
			b.Fatal("parse failed")
		}
	}
}

// ============================================================================
// Headers Benchmarks
// ============================================================================

func BenchmarkHeaders_Set(b *testing.B) {
	h := http.NewHeaders()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Set("Content-Type", "application/json")
	}
}

func BenchmarkHeaders_Get(b *testing.B) {
	h := http.NewHeaders()
	h.Set("Content-Type", "application/json")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = h.Get("Content-Type")
	}
}

func BenchmarkHeaders_Clone(b *testing.B) {
	h := makeHeaders(20)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = h.Clone()
	}
}

// ============================================================================
// Body Reader/Writer Benchmarks
// ============================================================================

func BenchmarkBodyReader_ReadAll_1KB(b *testing.B) {
	data := make([]byte, 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(data)
		result := reader.ReadAll()
		if result.IsErr() {
			b.Fatal("read failed")
		}
	}
}

func BenchmarkBodyReader_ReadAll_100KB(b *testing.B) {
	data := make([]byte, 100*1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := http.NewBodyReaderFromBytes(data)
		result := reader.ReadAll()
		if result.IsErr() {
			b.Fatal("read failed")
		}
	}
}

func BenchmarkBodyWriter_Write_1KB(b *testing.B) {
	data := make([]byte, 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer := http.NewBodyWriter()
		result := writer.Write(data)
		if result.IsErr() {
			b.Fatal("write failed")
		}
	}
}

func BenchmarkBodyWriter_WriteString(b *testing.B) {
	text := strings.Repeat("x", 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer := http.NewBodyWriter()
		result := writer.WriteString(text)
		if result.IsErr() {
			b.Fatal("write failed")
		}
	}
}

func BenchmarkBodyWriter_Json(b *testing.B) {
	user := makeTestUser(1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writer := http.NewBodyWriter()
		result := writer.Json(user)
		if result.IsErr() {
			b.Fatal("write failed")
		}
	}
}

// ============================================================================
// IntoResponse Benchmarks
// ============================================================================

func BenchmarkIntoResponse_String(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := http.StringResponse("Hello, World!")
		_ = s.IntoResponse()
	}
}

func BenchmarkIntoResponse_Bytes(b *testing.B) {
	data := make([]byte, 1024)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br := http.BytesResponse(data)
		_ = br.IntoResponse()
	}
}

func BenchmarkIntoResponse_StatusText(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		st := http.StatusTextResponse(http.StatusCreated, "Resource created")
		_ = st.IntoResponse()
	}
}

func BenchmarkIntoResponse_Html(b *testing.B) {
	html := "<html><body><h1>Hello</h1></body></html>"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := http.HtmlBody(html)
		_ = h.IntoResponse()
	}
}
