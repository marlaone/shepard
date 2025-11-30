package http

import (
	"net/url"
	"strings"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
)

// URL represents a parsed URL
type URL struct {
	Scheme   string
	Host     string
	Port     string
	Path     string
	RawQuery string
	Fragment string
}

// ParseURL parses a raw URL string into a URL struct
func ParseURL(rawURL string) shepard.Result[URL, error] {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return shepard.Err[URL, error](err)
	}

	port := parsed.Port()
	host := parsed.Hostname()

	return shepard.Ok[URL, error](URL{
		Scheme:   parsed.Scheme,
		Host:     host,
		Port:     port,
		Path:     parsed.Path,
		RawQuery: parsed.RawQuery,
		Fragment: parsed.Fragment,
	})
}

// String returns the string representation of the URL
func (u URL) String() string {
	var b strings.Builder

	if u.Scheme != "" {
		b.WriteString(u.Scheme)
		b.WriteString("://")
	}

	if u.Host != "" {
		b.WriteString(u.Host)
		if u.Port != "" {
			b.WriteString(":")
			b.WriteString(u.Port)
		}
	}

	b.WriteString(u.Path)

	if u.RawQuery != "" {
		b.WriteString("?")
		b.WriteString(u.RawQuery)
	}

	if u.Fragment != "" {
		b.WriteString("#")
		b.WriteString(u.Fragment)
	}

	return b.String()
}

// Query parses the RawQuery and returns QueryValues
func (u URL) Query() shepard.Result[QueryValues, error] {
	return ParseQuery(u.RawQuery)
}

// QueryValues stores URL query parameters
type QueryValues struct {
	values hashmap.HashMap[string, slice.Slice[string]]
}

// NewQueryValues creates a new empty QueryValues
func NewQueryValues() QueryValues {
	return QueryValues{
		values: hashmap.New[string, slice.Slice[string]](),
	}
}

// Set sets a query parameter value, replacing any existing values
func (qv *QueryValues) Set(key, value string) {
	qv.values.Insert(key, slice.Init(value))
}

// Add adds a query parameter value, keeping existing values
func (qv *QueryValues) Add(key, value string) {
	entry := qv.values.Entry(key)
	if entry.IsOccupied() {
		entry.AndModify(func(values *slice.Slice[string]) {
			values.Push(value)
		})
	} else {
		entry.OrInsert(slice.Init(value))
	}
}

// Get returns the first value for the given key
func (qv *QueryValues) Get(key string) shepard.Option[*string] {
	valuesOpt := qv.values.Get(key)
	if valuesOpt.IsNone() {
		return shepard.None[*string]()
	}
	return valuesOpt.Unwrap().First()
}

// GetAll returns all values for the given key
func (qv *QueryValues) GetAll(key string) slice.Slice[string] {
	valuesOpt := qv.values.Get(key)
	if valuesOpt.IsNone() {
		return slice.New[string]()
	}
	return *valuesOpt.Unwrap()
}

// Has returns true if the key exists
func (qv *QueryValues) Has(key string) bool {
	return qv.values.ContainsKey(key)
}

// Remove removes a key and all its values
func (qv *QueryValues) Remove(key string) {
	qv.values.Remove(key)
}

// Encode encodes the query values as a URL-encoded string
func (qv *QueryValues) Encode() string {
	stdValues := url.Values{}
	qv.values.Iter().Foreach(func(_ int, pair hashmap.Pair[*string, *slice.Slice[string]]) {
		pair.Value.Iter().Foreach(func(_ int, v string) {
			stdValues.Add(*pair.Key, v)
		})
	})
	return stdValues.Encode()
}

// ParseQuery parses a URL query string into QueryValues
func ParseQuery(query string) shepard.Result[QueryValues, error] {
	parsed, err := url.ParseQuery(query)
	if err != nil {
		return shepard.Err[QueryValues, error](err)
	}

	qv := NewQueryValues()
	for key, values := range parsed {
		for _, v := range values {
			qv.Add(key, v)
		}
	}

	return shepard.Ok[QueryValues, error](qv)
}
