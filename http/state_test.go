package http

import (
	"testing"

	"github.com/marlaone/shepard"
)

// Test types for typed state
type AppConfig struct {
	Debug bool
	Port  int
}

type DatabaseConn struct {
	Host string
}

func TestSetTypedState_GetTypedState(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	config := AppConfig{Debug: true, Port: 8080}
	SetTypedState(ctx, config)

	result := GetTypedState[AppConfig](ctx)
	if result.IsNone() {
		t.Fatal("expected Some, got None")
	}

	retrieved := result.Unwrap()
	if retrieved.Debug != true {
		t.Error("expected Debug=true")
	}
	if retrieved.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", retrieved.Port)
	}
}

func TestGetTypedState_NotFound(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	result := GetTypedState[AppConfig](ctx)
	if result.IsSome() {
		t.Error("expected None for missing state")
	}
}

func TestSetTypedState_MultipleTypes(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	config := AppConfig{Debug: true, Port: 8080}
	db := DatabaseConn{Host: "localhost"}

	SetTypedState(ctx, config)
	SetTypedState(ctx, db)

	// Both should be retrievable
	configResult := GetTypedState[AppConfig](ctx)
	if configResult.IsNone() {
		t.Fatal("expected AppConfig to be present")
	}
	if configResult.Unwrap().Port != 8080 {
		t.Error("AppConfig was overwritten")
	}

	dbResult := GetTypedState[DatabaseConn](ctx)
	if dbResult.IsNone() {
		t.Fatal("expected DatabaseConn to be present")
	}
	if dbResult.Unwrap().Host != "localhost" {
		t.Error("DatabaseConn has wrong value")
	}
}

func TestSetTypedState_Overwrite(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	SetTypedState(ctx, AppConfig{Debug: false, Port: 3000})
	SetTypedState(ctx, AppConfig{Debug: true, Port: 8080})

	result := GetTypedState[AppConfig](ctx)
	if result.IsNone() {
		t.Fatal("expected Some")
	}

	config := result.Unwrap()
	if config.Debug != true || config.Port != 8080 {
		t.Error("expected second value to overwrite first")
	}
}

func TestMustTypedState(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	config := AppConfig{Debug: true, Port: 8080}
	SetTypedState(ctx, config)

	retrieved := MustTypedState[AppConfig](ctx)
	if retrieved.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", retrieved.Port)
	}
}

func TestMustTypedState_Panics(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing state")
		}
	}()

	MustTypedState[AppConfig](ctx)
}

func TestTypedStateOr(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	defaultConfig := AppConfig{Debug: false, Port: 3000}

	// Not set - should return default
	result := TypedStateOr(ctx, defaultConfig)
	if result.Port != 3000 {
		t.Errorf("expected default Port=3000, got %d", result.Port)
	}

	// Set value - should return actual
	SetTypedState(ctx, AppConfig{Debug: true, Port: 8080})
	result = TypedStateOr(ctx, defaultConfig)
	if result.Port != 8080 {
		t.Errorf("expected actual Port=8080, got %d", result.Port)
	}
}

func TestHasTypedState(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	if HasTypedState[AppConfig](ctx) {
		t.Error("expected false for missing state")
	}

	SetTypedState(ctx, AppConfig{})
	if !HasTypedState[AppConfig](ctx) {
		t.Error("expected true after setting state")
	}
}

func TestRemoveTypedState(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	SetTypedState(ctx, AppConfig{Debug: true, Port: 8080})

	if !HasTypedState[AppConfig](ctx) {
		t.Fatal("state should exist before removal")
	}

	RemoveTypedState[AppConfig](ctx)

	if HasTypedState[AppConfig](ctx) {
		t.Error("state should not exist after removal")
	}
}

func TestWithState_Middleware(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	config := AppConfig{Debug: true, Port: 8080}
	middleware := WithState(config)

	var capturedConfig AppConfig
	handler := func() shepard.Result[Response, HttpError] {
		capturedConfig = MustTypedState[AppConfig](ctx)
		return TextOk("ok")
	}

	result := middleware(ctx, handler)
	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	if capturedConfig.Port != 8080 {
		t.Errorf("expected Port=8080 in handler, got %d", capturedConfig.Port)
	}
}

// Named State Tests

func TestSetNamedState_GetNamedState(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	SetNamedState(ctx, "primary", DatabaseConn{Host: "primary.db"})
	SetNamedState(ctx, "replica", DatabaseConn{Host: "replica.db"})

	primary := GetNamedState[DatabaseConn](ctx, "primary")
	if primary.IsNone() {
		t.Fatal("expected primary db")
	}
	if primary.Unwrap().Host != "primary.db" {
		t.Error("wrong primary host")
	}

	replica := GetNamedState[DatabaseConn](ctx, "replica")
	if replica.IsNone() {
		t.Fatal("expected replica db")
	}
	if replica.Unwrap().Host != "replica.db" {
		t.Error("wrong replica host")
	}
}

func TestGetNamedState_NotFound(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	result := GetNamedState[DatabaseConn](ctx, "missing")
	if result.IsSome() {
		t.Error("expected None for missing named state")
	}
}

func TestMustNamedState(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	SetNamedState(ctx, "mydb", DatabaseConn{Host: "localhost"})

	db := MustNamedState[DatabaseConn](ctx, "mydb")
	if db.Host != "localhost" {
		t.Error("wrong host")
	}
}

func TestMustNamedState_Panics(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for missing named state")
		}
	}()

	MustNamedState[DatabaseConn](ctx, "missing")
}

func TestNamedStateOr(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	defaultDB := DatabaseConn{Host: "default.db"}

	result := NamedStateOr(ctx, "missing", defaultDB)
	if result.Host != "default.db" {
		t.Error("expected default value")
	}

	SetNamedState(ctx, "actual", DatabaseConn{Host: "actual.db"})
	result = NamedStateOr(ctx, "actual", defaultDB)
	if result.Host != "actual.db" {
		t.Error("expected actual value")
	}
}

func TestWithNamedState_Middleware(t *testing.T) {
	req := NewRequest(MethodGet, "/test", EmptyBodyReader())
	ctx := NewContext(req)

	middleware := WithNamedState("primary", DatabaseConn{Host: "primary.db"})

	var capturedDB DatabaseConn
	handler := func() shepard.Result[Response, HttpError] {
		capturedDB = MustNamedState[DatabaseConn](ctx, "primary")
		return TextOk("ok")
	}

	result := middleware(ctx, handler)
	if result.IsErr() {
		t.Fatal("expected Ok result")
	}

	if capturedDB.Host != "primary.db" {
		t.Error("expected primary.db in handler")
	}
}
