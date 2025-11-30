package http

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/collections/slice"
)

type Userinfo struct {
	Username string
	Password string
}

type Values struct {
	values hashmap.HashMap[string, slice.Slice[string]]
}

func NewValues() *Values {
	return &Values{
		values: hashmap.New[string, slice.Slice[string]](),
	}
}

func (v *Values) Add(key, value string) {
	v.values.Entry(key).AndModify(func(s *slice.Slice[string]) {
		s.Push(value)
	}).OrInsert(slice.Init(value))
}

func (v *Values) Del(key string) {
	v.values.Remove(key)
}

func (v *Values) Get(key string) slice.Slice[string] {
	val := v.values.Get(key)
	if val.IsNone() {
		return slice.New[string]()
	}
	return *val.Unwrap()
}

func (v *Values) Has(key string) bool {
	return v.values.ContainsKey(key)
}

func (v *Values) Set(key, value string) {
	v.values.Entry(key).AndModify(func(s *slice.Slice[string]) {
		s.Clear()
		s.Push(value)
	}).OrInsert(slice.Init(value))
}

func ParseQuery(query string) shepard.Result[Values, error] {
	values := NewValues()

	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.Contains(key, ";") {
			return shepard.Err[Values, error](errors.New("invalid semicolon in query"))
		}

		if key == "" {
			continue
		}

		key, value, _ := strings.Cut(key, "=")
		key, err := url.QueryUnescape(key)
		if key == "" {
			continue
		}
		if err != nil {
			return shepard.Err[Values, error](err)
		}
		value, err = url.QueryUnescape(value)
		if err != nil {
			return shepard.Err[Values, error](err)
		}

		values.Add(key, strings.TrimSpace(value))
	}
	return shepard.Ok[Values, error](*values)
}

type URL struct {
	Scheme   string
	User     *Userinfo
	Host     string
	Port     uint16
	Path     string
	RawQuery string
	Fragment string

	values *Values
}

func (u URL) Default() URL {
	return URL{}
}

func (u *URL) Query() *Values {
	if u.values == nil {
		values := ParseQuery(u.RawQuery)
		if values.IsErr() {
			return NewValues()
		}
		u.values = new(Values)
		*u.values = values.Unwrap()
	}
	return u.values
}

func ParseURL(rawURL string) shepard.Result[URL, error] {

	u, err := url.Parse(rawURL)
	if err != nil {
		return shepard.Err[URL, error](err)
	}

	var port uint16

	if u.Port() != "" {
		p, err := strconv.ParseUint(u.Port(), 10, 16)
		if err != nil {
			return shepard.Err[URL, error](errors.New("invalid port"))
		}
		port = uint16(p)
	} else {
		switch u.Scheme {
		case "http":
			port = 80
		case "https":
			port = 443
		}
	}

	password, _ := u.User.Password()

	var userinfo *Userinfo

	if u.User.Username() != "" {
		userinfo = &Userinfo{
			Username: u.User.Username(),
			Password: password,
		}
	}

	return shepard.Ok[URL, error](URL{
		Scheme:   u.Scheme,
		User:     userinfo,
		Host:     u.Hostname(),
		Port:     port,
		Path:     u.Path,
		RawQuery: u.RawQuery,
		Fragment: u.Fragment,
	})
}

func ParseRequestURI(rawURI string) shepard.Result[URL, error] {
	u, err := url.ParseRequestURI(rawURI)
	if err != nil {
		return shepard.Err[URL, error](err)
	}

	return shepard.Ok[URL, error](URL{
		Path:     u.Path,
		RawQuery: u.RawQuery,
		Fragment: u.Fragment,
	})
}
