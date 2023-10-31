package http

type HttpResponseBytes struct {
	body Body

	statusCode StatusCode
	version    Version
	headers    *Headers
}

var _ Response[Body] = &HttpResponseBytes{}

func NewHttpResponseBytes() *HttpResponseBytes {
	headers := Headers{}.Default()
	return &HttpResponseBytes{
		statusCode: StatusCodeOk,
		version:    "1.1",
		headers:    &headers,
		body:       NewBytesBody(),
	}
}

func (r *HttpResponseBytes) SetStatusCode(statusCode StatusCode) {
	r.statusCode = statusCode
}

func (r *HttpResponseBytes) StatusCode() StatusCode {
	return r.statusCode
}

func (r *HttpResponseBytes) SetVersion(version Version) {
	r.version = version
}

func (r *HttpResponseBytes) Version() Version {
	return r.version
}

func (r *HttpResponseBytes) SetHeader(key string, values ...string) {
	r.headers.Set(key, values...)
}

func (r *HttpResponseBytes) Headers() *Headers {
	return r.headers
}

func (r *HttpResponseBytes) SetBody(body Body) {
	r.body = body
}

func (r *HttpResponseBytes) Finish() {
	r.body.Finish()
}

func (r *HttpResponseBytes) Body() Body {
	return r.body
}
