package http

import (
	"strings"
	"sync"

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
	values      []string // Native slice for performance
}

// Headers stores HTTP headers with case-insensitive keys
// Uses native Go map internally for fast parsing, with lazy HashMap conversion for Entry API
type Headers struct {
	// Native map for hot path (parsing, Get, Set, Append)
	// Stores headerValue by value to avoid heap allocation per header
	data map[string]headerValue

	// Lazy-loaded HashMap for Entry API access
	hashmap     *hashmap.HashMap[string, headerValue]
	hashmapOnce sync.Once
}

// NewHeaders creates a new empty Headers
func NewHeaders() *Headers {
	return &Headers{
		data: make(map[string]headerValue, 8), // Pre-allocate for typical header count
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
	h.data[canonical] = headerValue{
		originalKey: key,
		values:      []string{value},
	}
}

// Get returns the first value for the given header key
// Returns Option[string] - Some if found, None if not
func (h *Headers) Get(key string) shepard.Option[string] {
	canonical := canonicalKey(key)
	if hv, ok := h.data[canonical]; ok && len(hv.values) > 0 {
		return shepard.Some(hv.values[0])
	}
	return shepard.None[string]()
}

// GetAll returns all values for the given header key
func (h *Headers) GetAll(key string) slice.Slice[string] {
	canonical := canonicalKey(key)
	if hv, ok := h.data[canonical]; ok {
		// Return a copy as slice.Slice
		return slice.Init(hv.values...)
	}
	return slice.New[string]()
}

// Append adds a value to the header, keeping existing values
func (h *Headers) Append(key, value string) {
	canonical := canonicalKey(key)
	if hv, ok := h.data[canonical]; ok {
		hv.values = append(hv.values, value)
		h.data[canonical] = hv // Reassign since map values aren't addressable
	} else {
		h.data[canonical] = headerValue{
			originalKey: key,
			values:      []string{value},
		}
	}
}

// Remove removes a header by key
func (h *Headers) Remove(key string) {
	canonical := canonicalKey(key)
	delete(h.data, canonical)
}

// Contains returns true if the header exists
func (h *Headers) Contains(key string) bool {
	canonical := canonicalKey(key)
	_, ok := h.data[canonical]
	return ok
}

// Len returns the number of unique header keys
func (h *Headers) Len() int {
	return len(h.data)
}

// Clone creates a deep copy of the headers
func (h *Headers) Clone() *Headers {
	cloned := &Headers{
		data: make(map[string]headerValue, len(h.data)),
	}
	for k, v := range h.data {
		// Deep copy the values slice
		valuesCopy := make([]string, len(v.values))
		copy(valuesCopy, v.values)
		cloned.data[k] = headerValue{
			originalKey: v.originalKey,
			values:      valuesCopy,
		}
	}
	return cloned
}

// Iter returns an iterator over all header entries
func (h *Headers) Iter() iter.Iter[HeaderEntry] {
	entries := make([]HeaderEntry, 0, len(h.data))
	for _, hv := range h.data {
		for _, value := range hv.values {
			entries = append(entries, HeaderEntry{
				Key:   hv.originalKey,
				Value: value,
			})
		}
	}
	return iter.New(entries)
}

// Foreach iterates over all header keys and their values
func (h *Headers) Foreach(fn func(key string, values []string)) {
	for _, hv := range h.data {
		fn(hv.originalKey, hv.values)
	}
}

// Entry provides HashMap Entry API access for advanced use cases.
// This lazily initializes a HashMap from the native map data.
// Note: Mutations through Entry will modify the HashMap, not the native map.
// After calling Entry(), the Headers switches to HashMap mode for consistency.
func (h *Headers) Entry(key string) hashmap.Entry[string, headerValue] {
	h.ensureHashMap()
	return h.hashmap.Entry(canonicalKey(key))
}

// ensureHashMap lazily initializes the HashMap from native map data
func (h *Headers) ensureHashMap() {
	h.hashmapOnce.Do(func() {
		hm := hashmap.WithCapacity[string, headerValue](len(h.data))
		for k, v := range h.data {
			hm.Insert(k, v)
		}
		h.hashmap = &hm
	})
}

// AsHashMap returns the headers as a HashMap for advanced operations.
// This lazily initializes the HashMap if not already done.
func (h *Headers) AsHashMap() *hashmap.HashMap[string, headerValue] {
	h.ensureHashMap()
	return h.hashmap
}
