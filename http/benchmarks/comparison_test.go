// Package benchmarks provides comparative benchmarks between Shepard HTTP and net/http
package benchmarks

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	gohttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/marlaone/shepard/http"
)

// ============================================================================
// Shepard vs net/http Comparison Benchmarks
// ============================================================================

// --- Simple Text Response ---

func BenchmarkComparison_SimpleGet_NetHTTP(b *testing.B) {
	server := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	client := server.Client()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get(server.URL + "/")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkComparison_SimpleGet_Shepard(b *testing.B) {
	router := http.NewRouter()
	router.Get("/", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.Text("Hello, World!")
	}))

	// Start Shepard server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	addr := listener.Addr().String()
	listener.Close()

	server := http.NewServer(addr, router)
	go server.Serve()
	defer server.Stop()

	// Wait for server to start
	waitForServer(addr)

	client := &gohttp.Client{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get("http://" + addr + "/")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// --- JSON Response ---

func BenchmarkComparison_JSONResponse_NetHTTP(b *testing.B) {
	user := testUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30, Active: true}

	server := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	client := server.Client()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get(server.URL + "/json")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkComparison_JSONResponse_Shepard(b *testing.B) {
	user := testUser{ID: 1, Name: "John", Email: "john@example.com", Age: 30, Active: true}

	router := http.NewRouter()
	router.Get("/json", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		return http.JsonBody(user).IntoResponse()
	}))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	addr := listener.Addr().String()
	listener.Close()

	server := http.NewServer(addr, router)
	go server.Serve()
	defer server.Stop()

	waitForServer(addr)

	client := &gohttp.Client{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get("http://" + addr + "/json")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// --- Path Parameters ---

func BenchmarkComparison_PathParams_NetHTTP(b *testing.B) {
	mux := gohttp.NewServeMux()
	mux.HandleFunc("/users/{id}", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		id := r.PathValue("id")
		w.Write([]byte("User: " + id))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get(server.URL + "/users/123")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkComparison_PathParams_Shepard(b *testing.B) {
	router := http.NewRouter()
	router.Get("/users/:id", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		idOpt := ctx.Param("id")
		if idOpt.IsNone() {
			return http.TextWithStatus(http.StatusBadRequest, "missing id")
		}
		return http.Text("User: " + idOpt.Unwrap())
	}))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	addr := listener.Addr().String()
	listener.Close()

	server := http.NewServer(addr, router)
	go server.Serve()
	defer server.Stop()

	waitForServer(addr)

	client := &gohttp.Client{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get("http://" + addr + "/users/123")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// --- POST Simple (no body parsing) ---

func BenchmarkComparison_PostSimple_NetHTTP(b *testing.B) {
	server := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		io.Copy(io.Discard, r.Body)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	jsonData := `{"id":1,"name":"John","email":"john@example.com","age":30,"active":true}`
	client := server.Client()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Post(server.URL+"/echo", "application/json", strings.NewReader(jsonData))
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkComparison_PostSimple_Shepard(b *testing.B) {
	router := http.NewRouter()
	router.Post("/echo", http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
		// Just acknowledge the POST without parsing
		return http.Text("ok")
	}))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	addr := listener.Addr().String()
	listener.Close()

	server := http.NewServer(addr, router)
	go server.Serve()
	defer server.Stop()

	waitForServer(addr)

	jsonData := `{"id":1,"name":"John","email":"john@example.com","age":30,"active":true}`
	client := &gohttp.Client{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Post("http://"+addr+"/echo", "application/json", strings.NewReader(jsonData))
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// --- Multiple Routes (Router Matching) ---

func BenchmarkComparison_ManyRoutes_NetHTTP(b *testing.B) {
	mux := gohttp.NewServeMux()
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/api/v1/resource%d", i)
		mux.HandleFunc(path, func(w gohttp.ResponseWriter, r *gohttp.Request) {
			w.Write([]byte("ok"))
		})
	}

	server := httptest.NewServer(mux)
	defer server.Close()

	client := server.Client()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get(server.URL + "/api/v1/resource50")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkComparison_ManyRoutes_Shepard(b *testing.B) {
	router := http.NewRouter()
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/api/v1/resource%d", i)
		router.Get(path, http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
			return http.Text("ok")
		}))
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	addr := listener.Addr().String()
	listener.Close()

	server := http.NewServer(addr, router)
	go server.Serve()
	defer server.Stop()

	waitForServer(addr)

	client := &gohttp.Client{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get("http://" + addr + "/api/v1/resource50")
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

// --- Helper Functions ---

func waitForServer(addr string) {
	for i := 0; i < 100; i++ {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return
		}
	}
}

// ============================================================================
// Handler-Only Comparisons (No TCP overhead)
// These compare just the routing + handler execution layer
// ============================================================================

func BenchmarkHandlerOnly_SimpleGet_NetHTTP(b *testing.B) {
	handler := gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		w.Write([]byte("Hello, World!"))
	})

	req := httptest.NewRequest("GET", "/", nil)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

func BenchmarkHandlerOnly_SimpleGet_Shepard(b *testing.B) {
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

func BenchmarkHandlerOnly_PathParams_NetHTTP(b *testing.B) {
	mux := gohttp.NewServeMux()
	mux.HandleFunc("/users/{id}", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		id := r.PathValue("id")
		w.Write([]byte("User: " + id))
	})

	req := httptest.NewRequest("GET", "/users/123", nil)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkHandlerOnly_PathParams_Shepard(b *testing.B) {
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

func BenchmarkHandlerOnly_ManyRoutes_NetHTTP(b *testing.B) {
	mux := gohttp.NewServeMux()
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/api/v1/resource%d", i)
		mux.HandleFunc(path, func(w gohttp.ResponseWriter, r *gohttp.Request) {
			w.Write([]byte("ok"))
		})
	}

	req := httptest.NewRequest("GET", "/api/v1/resource50", nil)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkHandlerOnly_ManyRoutes_Shepard(b *testing.B) {
	router := http.NewRouter()
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/api/v1/resource%d", i)
		router.Get(path, http.SimpleHandlerFunc(func(ctx *http.Context) http.Response {
			return http.Text("ok")
		}))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := http.NewRequest(http.MethodGet, "/api/v1/resource50", http.EmptyBodyReader())
		ctx := http.NewContext(req)
		_ = router.Handle(ctx)
	}
}
