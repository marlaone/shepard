package http

import (
	"strings"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
	"github.com/marlaone/shepard/iter"
)

// Common header names
const (
	HeaderContentType      = "Content-Type"
	HeaderContentLength    = "Content-Length"
	HeaderAccept           = "Accept"
	HeaderAuthorization    = "Authorization"
	HeaderHost             = "Host"
	HeaderUserAgent        = "User-Agent"
	HeaderCacheControl     = "Cache-Control"
	HeaderConnection       = "Connection"
	HeaderCookie           = "Cookie"
	HeaderSetCookie        = "Set-Cookie"
	HeaderLocation         = "Location"
	HeaderContentEncoding  = "Content-Encoding"
	HeaderTransferEncoding = "Transfer-Encoding"
)

// canonicalHeaderCache pre-computes lowercase versions of common headers
// to avoid allocations from strings.ToLower on frequently used header names
var canonicalHeaderCache = map[string]string{
	// Content-Type variations
	"Content-Type": "content-type",
	"content-type": "content-type",
	"CONTENT-TYPE": "content-type",
	"Content-type": "content-type",
	// Content-Length variations
	"Content-Length": "content-length",
	"content-length": "content-length",
	"CONTENT-LENGTH": "content-length",
	"Content-length": "content-length",
	// Accept variations
	"Accept": "accept",
	"accept": "accept",
	"ACCEPT": "accept",
	// Authorization variations
	"Authorization": "authorization",
	"authorization": "authorization",
	"AUTHORIZATION": "authorization",
	// Host variations
	"Host": "host",
	"host": "host",
	"HOST": "host",
	// User-Agent variations
	"User-Agent": "user-agent",
	"user-agent": "user-agent",
	"USER-AGENT": "user-agent",
	"User-agent": "user-agent",
	// Cache-Control variations
	"Cache-Control": "cache-control",
	"cache-control": "cache-control",
	"CACHE-CONTROL": "cache-control",
	"Cache-control": "cache-control",
	// Connection variations
	"Connection": "connection",
	"connection": "connection",
	"CONNECTION": "connection",
	// Cookie variations
	"Cookie": "cookie",
	"cookie": "cookie",
	"COOKIE": "cookie",
	// Set-Cookie variations
	"Set-Cookie": "set-cookie",
	"set-cookie": "set-cookie",
	"SET-COOKIE": "set-cookie",
	"Set-cookie": "set-cookie",
	// Location variations
	"Location": "location",
	"location": "location",
	"LOCATION": "location",
	// Content-Encoding variations
	"Content-Encoding": "content-encoding",
	"content-encoding": "content-encoding",
	"CONTENT-ENCODING": "content-encoding",
	"Content-encoding": "content-encoding",
	// Transfer-Encoding variations
	"Transfer-Encoding": "transfer-encoding",
	"transfer-encoding": "transfer-encoding",
	"TRANSFER-ENCODING": "transfer-encoding",
	"Transfer-encoding": "transfer-encoding",
	// Accept-Encoding variations
	"Accept-Encoding": "accept-encoding",
	"accept-encoding": "accept-encoding",
	"ACCEPT-ENCODING": "accept-encoding",
	"Accept-encoding": "accept-encoding",
	// Accept-Language variations
	"Accept-Language": "accept-language",
	"accept-language": "accept-language",
	"ACCEPT-LANGUAGE": "accept-language",
	// X-Requested-With variations
	"X-Requested-With": "x-requested-with",
	"x-requested-with": "x-requested-with",
	// X-Forwarded-For variations
	"X-Forwarded-For": "x-forwarded-for",
	"x-forwarded-for": "x-forwarded-for",
	"X-FORWARDED-FOR": "x-forwarded-for",
	// Origin variations
	"Origin": "origin",
	"origin": "origin",
	"ORIGIN": "origin",
	// Referer variations
	"Referer": "referer",
	"referer": "referer",
	"REFERER": "referer",
}

// HeaderEntry represents a single header key-value pair
type HeaderEntry struct {
	Key   string
	Value string
}

// headerValue stores the original key case and values
type headerValue struct {
	originalKey string
	values      slice.Slice[string]
}

// Headers stores HTTP headers with case-insensitive keys
// Uses HashMap internally for efficient lookups
type Headers struct {
	headers hashmap.HashMap[string, headerValue]
}

// NewHeaders creates a new empty Headers
func NewHeaders() *Headers {
	return &Headers{
		headers: hashmap.New[string, headerValue](),
	}
}

// canonicalKey returns the canonical form of the header key (lowercase)
// Uses cache lookup for common headers to avoid allocation
func canonicalKey(key string) string {
	if cached, ok := canonicalHeaderCache[key]; ok {
		return cached
	}
	return strings.ToLower(key)
}

// Set sets a header value, replacing any existing values
func (h *Headers) Set(key, value string) {
	canonical := canonicalKey(key)
	hv := headerValue{
		originalKey: key,
		values:      slice.Init(value),
	}
	h.headers.Insert(canonical, hv)
}

// Get returns the first value for the given header key
// Returns Option[*string] - Some if found, None if not
func (h *Headers) Get(key string) shepard.Option[*string] {
	canonical := canonicalKey(key)
	hvOpt := h.headers.Get(canonical)
	if hvOpt.IsNone() {
		return shepard.None[*string]()
	}
	hv := hvOpt.Unwrap()
	return hv.values.First()
}

// GetAll returns all values for the given header key
func (h *Headers) GetAll(key string) slice.Slice[string] {
	canonical := canonicalKey(key)
	hvOpt := h.headers.Get(canonical)
	if hvOpt.IsNone() {
		return slice.New[string]()
	}
	return hvOpt.Unwrap().values
}

// Append adds a value to the header, keeping existing values
func (h *Headers) Append(key, value string) {
	canonical := canonicalKey(key)
	entry := h.headers.Entry(canonical)
	if entry.IsOccupied() {
		entry.AndModify(func(hv *headerValue) {
			hv.values.Push(value)
		})
	} else {
		entry.OrInsert(headerValue{
			originalKey: key,
			values:      slice.Init(value),
		})
	}
}

// Remove removes a header by key
func (h *Headers) Remove(key string) {
	canonical := canonicalKey(key)
	h.headers.Remove(canonical)
}

// Contains returns true if the header exists
func (h *Headers) Contains(key string) bool {
	canonical := canonicalKey(key)
	return h.headers.ContainsKey(canonical)
}

// Len returns the number of unique header keys
func (h *Headers) Len() int {
	return h.headers.Len()
}

// Clone creates a deep copy of the headers
func (h *Headers) Clone() *Headers {
	cloned := NewHeaders()
	h.headers.Iter().Foreach(func(_ int, pair hashmap.Pair[*string, *headerValue]) {
		// Clone the values slice
		clonedValues := slice.WithCapacity[string](pair.Value.values.Len())
		pair.Value.values.Iter().Foreach(func(_ int, v string) {
			clonedValues.Push(v)
		})
		cloned.headers.Insert(*pair.Key, headerValue{
			originalKey: pair.Value.originalKey,
			values:      clonedValues,
		})
	})
	return cloned
}

// Iter returns an iterator over all header entries
func (h *Headers) Iter() iter.Iter[HeaderEntry] {
	entries := slice.New[HeaderEntry]()
	h.headers.Iter().Foreach(func(_ int, pair hashmap.Pair[*string, *headerValue]) {
		pair.Value.values.Iter().Foreach(func(_ int, value string) {
			entries.Push(HeaderEntry{
				Key:   pair.Value.originalKey,
				Value: value,
			})
		})
	})
	return entries.Iter()
}

// Foreach iterates over all header keys and their values
func (h *Headers) Foreach(fn func(key string, values []string)) {
	h.headers.Iter().Foreach(func(_ int, pair hashmap.Pair[*string, *headerValue]) {
		values := make([]string, pair.Value.values.Len())
		pair.Value.values.Iter().Foreach(func(i int, v string) {
			values[i] = v
		})
		fn(pair.Value.originalKey, values)
	})
}
