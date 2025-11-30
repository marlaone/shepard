package http

// StatusCode represents an HTTP status code
type StatusCode uint16

// 1xx Informational
const (
	StatusContinue           StatusCode = 100
	StatusSwitchingProtocols StatusCode = 101
	StatusProcessing         StatusCode = 102
	StatusEarlyHints         StatusCode = 103
)

// 2xx Success
const (
	StatusOK                   StatusCode = 200
	StatusCreated              StatusCode = 201
	StatusAccepted             StatusCode = 202
	StatusNonAuthoritativeInfo StatusCode = 203
	StatusNoContent            StatusCode = 204
	StatusResetContent         StatusCode = 205
	StatusPartialContent       StatusCode = 206
)

// 3xx Redirection
const (
	StatusMultipleChoices   StatusCode = 300
	StatusMovedPermanently  StatusCode = 301
	StatusFound             StatusCode = 302
	StatusSeeOther          StatusCode = 303
	StatusNotModified       StatusCode = 304
	StatusUseProxy          StatusCode = 305
	StatusTemporaryRedirect StatusCode = 307
	StatusPermanentRedirect StatusCode = 308
)

// 4xx Client Errors
const (
	StatusBadRequest                   StatusCode = 400
	StatusUnauthorized                 StatusCode = 401
	StatusPaymentRequired              StatusCode = 402
	StatusForbidden                    StatusCode = 403
	StatusNotFound                     StatusCode = 404
	StatusMethodNotAllowed             StatusCode = 405
	StatusNotAcceptable                StatusCode = 406
	StatusProxyAuthRequired            StatusCode = 407
	StatusRequestTimeout               StatusCode = 408
	StatusConflict                     StatusCode = 409
	StatusGone                         StatusCode = 410
	StatusLengthRequired               StatusCode = 411
	StatusPreconditionFailed           StatusCode = 412
	StatusRequestEntityTooLarge        StatusCode = 413
	StatusRequestURITooLong            StatusCode = 414
	StatusUnsupportedMediaType         StatusCode = 415
	StatusRequestedRangeNotSatisfiable StatusCode = 416
	StatusExpectationFailed            StatusCode = 417
	StatusTeapot                       StatusCode = 418
	StatusMisdirectedRequest           StatusCode = 421
	StatusUnprocessableEntity          StatusCode = 422
	StatusLocked                       StatusCode = 423
	StatusFailedDependency             StatusCode = 424
	StatusTooEarly                     StatusCode = 425
	StatusUpgradeRequired              StatusCode = 426
	StatusPreconditionRequired         StatusCode = 428
	StatusTooManyRequests              StatusCode = 429
	StatusRequestHeaderFieldsTooLarge  StatusCode = 431
	StatusUnavailableForLegalReasons   StatusCode = 451
)

// 5xx Server Errors
const (
	StatusInternalServerError           StatusCode = 500
	StatusNotImplemented                StatusCode = 501
	StatusBadGateway                    StatusCode = 502
	StatusServiceUnavailable            StatusCode = 503
	StatusGatewayTimeout                StatusCode = 504
	StatusHTTPVersionNotSupported       StatusCode = 505
	StatusVariantAlsoNegotiates         StatusCode = 506
	StatusInsufficientStorage           StatusCode = 507
	StatusLoopDetected                  StatusCode = 508
	StatusNotExtended                   StatusCode = 510
	StatusNetworkAuthenticationRequired StatusCode = 511
)

var statusTexts = map[StatusCode]string{
	StatusContinue:                      "Continue",
	StatusSwitchingProtocols:            "Switching Protocols",
	StatusProcessing:                    "Processing",
	StatusEarlyHints:                    "Early Hints",
	StatusOK:                            "OK",
	StatusCreated:                       "Created",
	StatusAccepted:                      "Accepted",
	StatusNonAuthoritativeInfo:          "Non-Authoritative Information",
	StatusNoContent:                     "No Content",
	StatusResetContent:                  "Reset Content",
	StatusPartialContent:                "Partial Content",
	StatusMultipleChoices:               "Multiple Choices",
	StatusMovedPermanently:              "Moved Permanently",
	StatusFound:                         "Found",
	StatusSeeOther:                      "See Other",
	StatusNotModified:                   "Not Modified",
	StatusUseProxy:                      "Use Proxy",
	StatusTemporaryRedirect:             "Temporary Redirect",
	StatusPermanentRedirect:             "Permanent Redirect",
	StatusBadRequest:                    "Bad Request",
	StatusUnauthorized:                  "Unauthorized",
	StatusPaymentRequired:               "Payment Required",
	StatusForbidden:                     "Forbidden",
	StatusNotFound:                      "Not Found",
	StatusMethodNotAllowed:              "Method Not Allowed",
	StatusNotAcceptable:                 "Not Acceptable",
	StatusProxyAuthRequired:             "Proxy Authentication Required",
	StatusRequestTimeout:                "Request Timeout",
	StatusConflict:                      "Conflict",
	StatusGone:                          "Gone",
	StatusLengthRequired:                "Length Required",
	StatusPreconditionFailed:            "Precondition Failed",
	StatusRequestEntityTooLarge:         "Request Entity Too Large",
	StatusRequestURITooLong:             "Request URI Too Long",
	StatusUnsupportedMediaType:          "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable:  "Requested Range Not Satisfiable",
	StatusExpectationFailed:             "Expectation Failed",
	StatusTeapot:                        "I'm a teapot",
	StatusMisdirectedRequest:            "Misdirected Request",
	StatusUnprocessableEntity:           "Unprocessable Entity",
	StatusLocked:                        "Locked",
	StatusFailedDependency:              "Failed Dependency",
	StatusTooEarly:                      "Too Early",
	StatusUpgradeRequired:               "Upgrade Required",
	StatusPreconditionRequired:          "Precondition Required",
	StatusTooManyRequests:               "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:   "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:    "Unavailable For Legal Reasons",
	StatusInternalServerError:           "Internal Server Error",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
	StatusInsufficientStorage:           "Insufficient Storage",
	StatusLoopDetected:                  "Loop Detected",
	StatusNotExtended:                   "Not Extended",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

// Int returns the integer value of the status code
func (s StatusCode) Int() int {
	return int(s)
}

// Text returns the text description of the status code
func (s StatusCode) Text() string {
	if text, ok := statusTexts[s]; ok {
		return text
	}
	return ""
}

// IsInformational returns true for 1xx status codes
func (s StatusCode) IsInformational() bool {
	return s >= 100 && s < 200
}

// IsSuccess returns true for 2xx status codes
func (s StatusCode) IsSuccess() bool {
	return s >= 200 && s < 300
}

// IsRedirect returns true for 3xx status codes
func (s StatusCode) IsRedirect() bool {
	return s >= 300 && s < 400
}

// IsClientError returns true for 4xx status codes
func (s StatusCode) IsClientError() bool {
	return s >= 400 && s < 500
}

// IsServerError returns true for 5xx status codes
func (s StatusCode) IsServerError() bool {
	return s >= 500 && s < 600
}

// IsError returns true for 4xx and 5xx status codes
func (s StatusCode) IsError() bool {
	return s.IsClientError() || s.IsServerError()
}
